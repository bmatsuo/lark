local go = require('go')

go.default_sources = {
    './cmd/...',
    './luamodules/...',
}

lark.default_task = 'build'

lark.task{'gen', function ()
    go.gen()
end}

lark.task{'build', function ()
    lark.run{'gen'}
    go.build{'./cmd/...'}
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
