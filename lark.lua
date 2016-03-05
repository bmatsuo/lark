-- the following line migrates from the v0.4.0 task API to v0.5.0 API
lark.task = require('lark.task')

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
local novendor =
    lark.exec{'glide', 'novendor', stdout='$', stderr='/dev/null', echo=false}
string.gsub(novendor, '(%S+)', function(p) table.insert(sources, p) end)
go.default_sources = sources

all = lark.task .. function()
    lark.run('gen')
    lark.run('test')
    lark.run('build')
end
