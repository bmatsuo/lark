lark.task{'generate', function ()
    lark.exec{'go', 'generate', './cmd/...'}
end}


lark.task{'build', function ()
    lark.run{'generate'}
    lark.exec{'go', 'build', './cmd/...'}
end}

lark.task{'install', function ()
    lark.run{'generate'}
    lark.exec{'go', 'install', './cmd/...'}
end}
