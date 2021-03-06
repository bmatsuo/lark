local task = require('lark.task')
local go = require('go')

local function collect(...)
    local args = {unpack(arg)}
    args.stdout = '$'
    args.echo = false
    local output = lark.exec(args)
    local ret = {}
    local insert = function(x) table.insert(ret, x) end
    string.gsub(output, '(%S+)', insert)
    return ret
end

init = task .. function()
    lark.exec{'glide', 'install'}
end

clean = task .. function()
    lark.exec{'rm', '-f', 'lark', 'docgen'}
end

gen = task .. function ()
    go.gen()
    lark.run('gen_docs')
end

gen_docs = task .. function ()
    lark.run('./cmd/docgen')
    lark.exec{'./docgen'}
end

build = task .. function ()
    lark.run('./cmd/lark')
end

build_all = task .. function()
    local cmds = collect('sh', '-c', 'ls -d ./cmd/*')
    for _, build in pairs(cmds) do
        lark.run(build)
    end
end

build_patt = task.pattern[[^./cmd/.*]] .. function (ctx)
    go.build{task.get_name(ctx)}
end

install = task .. function ()
    go.install{'./cmd/lark'}
end

test = task .. function(ctx)
    local race = task.get_param(ctx, 'race')
    if race then
        go.test{race=true}
    else
        go.test{cover=true}
    end
end
