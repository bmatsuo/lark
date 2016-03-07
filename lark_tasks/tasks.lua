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
    lark.run('./cmd/lark')
end

build_patt = task.pattern[[^./cmd/.*]] .. function (ctx)
    go.build{task.get_name(ctx)}
end

install = task .. function ()
    go.install()
end

test = task .. function(ctx)
    local race = task.get_param(ctx, 'race')
    if race then
        go.test{race=true}
    else
        go.test{cover=true}
    end
end
