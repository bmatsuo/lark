local task = require('lark.task')
local go = require('go')

init = task .. function()
    lark.exec{'glide', 'install'}
end

clean = task .. function()
    lark.exec{'rm', '-f', 'lark'}
end

gen = task .. function ()
    go.gen()
end

build = task .. function ()
    go.build{'./cmd/...', ldflags=ldflags}
end

install = task .. function ()
    go.install{ldflags=ldflags}
end

test = task .. function(ctx)
    local race = task.get_param(ctx, 'race')
    if race then
        go.test{race=true, ldflags=ldflags}
    else
        go.test{cover=true, ldflags=ldflags}
    end
end
