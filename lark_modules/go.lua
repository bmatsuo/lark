local doc = require('doc')
local go =
    doc.desc[[
    Functions to assist in running the go command and its subcommands.  Most
    functions in this module accept a table of options.  Some of these options
    are common to all go subcommands.  These options are the strings
    "asmflags", "gcflags", "ldflags", and the boolean "race".

        > go.build{gcflags="-B", race=true}
    ]] ..
    doc.sig[[(...) => ()]] ..
    doc.param[[...  string
    -- Source trees in which to generate code.  If values are present they will
    override the value of default_sources.  Any arrays given will be
    recursively flattened to produce a source tree list.
    ]] ..
    doc.var[[default_sources  The default list of source (trees) passed to go subcommands.]] ..
    {
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
    if val then
        if type(val) == 'boolean' then
            insert_args(tcmd, {flag})
        else
            insert_args(tcmd, {flag, val})
        end
    end
end

local function insert_common_build_flags(tcmd, opt)
    opt_flag(tcmd, '-asmflags', opt.asmflags)
    opt_flag(tcmd, '-gcflags', opt.gcflags)
    opt_flag(tcmd, '-ldflags', opt.ldflags)
    opt_flag(tcmd, '-race', opt.race)
end

local function flatten(...)
    local flat = {}
    for _, v in pairs(arg) do
        if type(v) ~= 'table' then
            table.insert(flat, v)
        else
            for _, v in pairs(flatten(unpack(v))) do
                table.insert(flat, v)
            end
        end
    end
    return flat
end

go.gen =
    doc.desc[[Run the ``go generate'' command.]] ..
    doc.param[[sources  table
    -- If array values are present in sources they will override the value of
    default_sources.
    ]] ..
    function(...)
        local opt = flatten(arg)
        local cmd = {'go', 'generate'}
        if #arg == 0 then
            insert_args(cmd, go.default_sources)
            lark.exec{cmd}
            return
        end

        local args = opt
        if #args == 0 then
            args = go.default_sources
        end
        insert_args(cmd, args)

        lark.exec{cmd}
    end

go.install =
    doc.desc[[Run the ``go install'' command.]] ..
    doc.param[[opt  table
    -- If array values are present in sources they will override the value of
    default_sources.  All common options are supported for the install()
    function. 
    ]] ..
    function(opt)
        local cmd = {'go', 'install'}
        if not opt then
            insert_args(cmd, go.default_sources)
            lark.exec{cmd}
            return
        end

        insert_common_build_flags(cmd, opt)

        local args = opt
        if #args == 0 then
            args = go.default_sources
        end
        insert_args(cmd, args)

        lark.exec{cmd}
    end

go.build =
    doc.desc[[Run the ``go build'' command.]] ..
    doc.param[[opt  table
    -- If array values are present in sources they will override the value of
    default_sources.  All common options are supported for the build()
    function. 
    ]] ..
    function(opt)
        local cmd = {'go', 'build'}
        if not opt then
            insert_args(cmd, go.default_sources)
            lark.exec{cmd}
            return
        end

        insert_common_build_flags(cmd, opt)

        local args = opt
        if #args == 0 then
            args = go.default_sources
        end
        insert_args(cmd, args)

        lark.exec{cmd}
    end

go.test =
    doc.desc[[Run the ``go test'' command.]] ..
    doc.param[[opt  table
    -- If array values are present in sources they will override the value of
    default_sources.  All common options are supported for the test()
    function. 
    ]] ..
    function(opt)
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
        if #args == 0 then
            args = go.default_sources
        end
        insert_args(cmd, args)

        lark.exec{cmd}
    end

return go
