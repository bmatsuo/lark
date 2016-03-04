local go = require('go')

lark.task('init', function()
    lark.exec{'glide', 'install'}
end)

lark.task('clean',  function()
    lark.exec{'rm', '-f', 'lark'}
end)

lark.task('gen', function ()
    go.gen()
end)

lark.task('build', function ()
    go.build{'./cmd/...', ldflags=ldflags}
end)

lark.task('install', function ()
    go.install{ldflags=ldflags}
end)

lark.task('test', function(ctx)
    local race = lark.get_param(ctx, 'race')
    if race then
        go.test{race=true}
    else
        go.test{cover=true}
    end
end)
