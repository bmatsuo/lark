local go = require('go')

local import = 'github.com/bmatsuo/lark'
local version = '0.1.0-dev'

local _ldflags = {
    string.format('-X %s/larkmeta.Version=%s', import, version),
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
    lark.run{'gen'}
    go.build{'./cmd/...', ldflags=ldflags}
end}

lark.task{'install', function ()
    lark.run{'gen'}
    go.install()
end}

-- BUG: We don't want to test the vendored packages.  But we want to run the
-- tests for everything else.
lark.task{'test', function()
    lark.run{'gen'}
    go.test{cover=true}
end}
