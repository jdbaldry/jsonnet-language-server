{ pkgs ? import <nixpkgs> }:

with pkgs;
mkShell {
  buildInputs = [
    go_1_16
    golangci-lint
    gopls
    snitch
  ];
  shellHook = ''
    export PATH="$(pwd)":"$PATH"
  '';
}
