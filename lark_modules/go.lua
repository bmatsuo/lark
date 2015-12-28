local go = {
    default_sources = {'.'}
}

local function insert_args(tcmd, targs)
    for i, arg in pairs(targs) do
        if type(i) == 'number' then
            table.insert(tcmd, arg)
        end
    end
end

-- this function has some serious problems. but whatever for now, it's local.
local function opt_flag(tcmd, flag, val)
    if val then insert_args(tcmd, {flag, val}) end
end

local function insert_common_build_flags(tcmd, opt)
    opt_flag(tcmd, '-asmflags', opt.asmflags)
    opt_flag(tcmd, '-gcflags', opt.gcflags)
    opt_flag(tcmd, '-ldflags', opt.ldflags)
end

go.gen = function(opt)
    local cmd = {'go', 'generate'}
    if not opt then
        insert_args(cmd, go.default_sources)
        lark.exec{cmd}
        return
    end

    local args = opt
    if table.getn(args) == 0 then
        args = go.default_sources
    end
    insert_args(cmd, args)

    lark.exec{cmd}
end

go.install = function(opt)
    local cmd = {'go', 'install'}
    if not opt then
        insert_args(cmd, go.default_sources)
        lark.exec{cmd}
        return
    end

    insert_common_build_flags(cmd, opt)

    local args = opt
    if table.getn(args) == 0 then
        args = go.default_sources
    end
    insert_args(cmd, args)

    lark.exec{cmd}
end

go.build = function(opt)
    local cmd = {'go', 'build'}
    if not opt then
        insert_args(cmd, go.default_sources)
        lark.exec{cmd}
        return
    end

    insert_common_build_flags(cmd, opt)

    local args = opt
    if table.getn(args) == 0 then
        args = go.default_sources
    end
    insert_args(cmd, args)

    lark.exec{cmd}
end

go.test = function(opt)
    local cmd = {'go', 'test'}
    if not opt then
        insert_args(cmd, go.default_sources)
        lark.exec{cmd}
        return
    end

    insert_common_build_flags(cmd, opt)

    if opt.cover then
        if type(opt.cover) == 'string' then
            local arg = string.format('-coverprofile=%s', opt.cover)
            table.insert(cmd, arg)
        else
            table.insert(cmd, '-cover')
        end
    end

    local args = opt
    if table.getn(args) == 0 then
        args = go.default_sources
    end
    insert_args(cmd, args)

    lark.exec{cmd}
end

return go
