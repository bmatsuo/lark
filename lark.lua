local path = require('path')
local go = require('go')
local version = require('version')

local import = 'github.com/bmatsuo/lark'

local _ldflags = {
    string.format('-X %s/larkmeta.Version=%s', import, version.get()),
}
local ldflags = table.concat(_ldflags, ' ')

go.default_sources = {
    './cmd/...',
    './luamodules/...',
}

lark.default_task = 'build'

lark.task{'clean', function()
    lark.exec{'rm', '-f', 'lark'}
end}

lark.task{'gen', function ()
    go.gen()
end}

lark.task{'build', function ()
    go.build{'./cmd/...', ldflags=ldflags}
end}

lark.task{'install', function ()
    go.install{ldflags=ldflags}
end}

-- BUG: We don't want to test the vendored packages.  But we want to run the
-- tests for everything else.
lark.task{'test', function()
    go.test{cover=true}
end}

lark.task{'release', function()
    local release_dir = 'release'
    local vx = version.get()
    local release_name_template = 'lark_' .. vx .. '_{{.OS}}_{{.Arch}}'
    local template = path.join(release_dir, release_name_template, '{{.Dir}}')
    lark.run{'gen', 'test'}
    lark.exec{'mkdir', '-p', 'release'}
    lark.exec{'gox', '-output='..template, './cmd/...'}
end}
