--local doc = require('doc')
local go = require('go')
local version = require('version')

-- configure Go ldflags
local import = 'github.com/bmatsuo/lark'
local _ldflags = {
    string.format('-X %s/larkmeta.Version=%s', import, version.get()),
}
ldflags = table.concat(_ldflags, ' ')

local sources = {}
local novendor = lark.exec{'glide', 'novendor', stdout='$', stderr='/dev/null'}
string.gsub(novendor, '(%S+)', function(p) table.insert(sources, p) end)
go.default_sources = sources

lark.default_task = 'all'

-- doc.help(doc.help)
