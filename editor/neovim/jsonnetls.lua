local util = require 'lspconfig/util'
local configs = require 'lspconfig/configs'

configs.jsonnetls = {
  default_config = {
    cmd = { 'jsonnet-language-server' };
    filetypes = { 'jsonnet', 'libsonnet' };
    root_dir = function(fname)
      return util.root_pattern('jsonnetfile.json')(fname) or util.find_git_ancestor(fname) or util.path.dirname(fname)
    end;
    settings = {};
    docs = {
      description = [[
https://github.com/jdbaldry/jsonnet-language-server/

A language server for the Jsonnet data templating language.
https://jsonnet.org/

`jsonnet-language-server` can be installed using `go`:

```sh
go install github.com/jdbaldry/jsonnet-language-server
```
      ]];
    }
  };
}
