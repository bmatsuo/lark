local go = require('go')
local version = require('version')

-- configure Go ldflags
local import = 'github.com/bmatsuo/lark'
local _ldflags = {
    string.format('-X %s/larkmeta.Version=%s', import, version.get()),
}
ldflags = table.concat(_ldflags, ' ')

go.default_sources = {
    './cmd/...',
    './luamodules/...',
}

lark.default_task = 'all'
