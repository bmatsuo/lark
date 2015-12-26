lark.task{'build', function ()
    lark.exec{'go', 'generate', './cmd/...'}
    lark.exec{'go', 'build', './cmd/...'}
end}
