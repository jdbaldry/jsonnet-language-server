local lspconfig = require 'lspconfig'

lspconfig.jsonnetls.setup{
  cmd = { '~/go/bin/jsonnet-language-server' };
}
