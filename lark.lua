lark.default_task = 'build'

lark.task{'gen', function ()
    lark.exec{'go', 'generate', './cmd/...'}
end}

lark.task{'build', function ()
    lark.run{'gen'}
    lark.exec{'go', 'build', './cmd/...'}
end}

lark.task{'install', function ()
    lark.run{'gen'}
    lark.exec{'go', 'install', './cmd/...'}
end}

-- BUG: We don't want to test the vendored packages.  But we want to run the
-- tests for everything else.
lark.task{'test', function()
    lark.run{'gen'}
    lark.exec{'go', 'test', '-cover', './cmd/...'}
end}
