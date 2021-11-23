// jsonnet-language-server: A Language Server Protocol server for Jsonnet.
// Copyright (C) 2021 Jack Baldry

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jdbaldry/go-language-server-protocol/jsonrpc2"
	"github.com/jdbaldry/go-language-server-protocol/lsp/protocol"
)

const (
	name = "jsonnet-language-server"
)

var (
	// Set with `-ldflags="-X 'main.version=<version>'"`
	version = "dev"
)

// printVersion prints version text to the provided writer.
func printVersion() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "%s version %s\n", name, version)
}

// printVersion prints help text to the provided writer.
func printHelp() {
	w := flag.CommandLine.Output()

	printVersion()
	fmt.Fprintln(w)
	flag.PrintDefaults()
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Environment variables:\n")
	fmt.Fprintf(w, "  JSONNET_PATH is a %q separated list of directories\n", filepath.ListSeparator)
	fmt.Fprintf(w, "  added in reverse order before the paths specified by --jpath.\n")
	fmt.Fprintf(w, "  These are equivalent:\n")
	fmt.Fprintf(w, "    JSONNET_PATH=a:b %s -J c -J d\n", name)
	fmt.Fprintf(w, "    JSONNET_PATH=d:c:a:b %s\n", name)
	fmt.Fprintf(w, "    %s -J b -J a -J c -J d\n", name)
}

func main() {
	jpaths := filepath.SplitList(os.Getenv("JSONNET_PATH"))

	flag.Usage = printHelp
	versionArg := flag.Bool("version", false, "Print version and exit")
	jpathArg := flag.String("jpath", "", "Specify an additional library search dir (right-most wins)")
	flag.Parse()

	if *versionArg {
		printVersion()
		os.Exit(0)
	}

	jpaths = append([]string{*jpathArg}, jpaths...)

	log.Println("Starting the language server")

	ctx := context.Background()
	stream := jsonrpc2.NewHeaderStream(stdio{})
	conn := jsonrpc2.NewConn(stream)
	client := protocol.ClientDispatcher(conn)

	s, err := newServer(client, jpaths)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	conn.Go(ctx, protocol.Handlers(
		protocol.ServerHandler(s, jsonrpc2.MethodNotFound)))
	<-conn.Done()
	if err := conn.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
