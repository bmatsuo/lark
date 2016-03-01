package doc

// DocLib contains Lua source code for the doc module.
var DocLib = `local signatures = setmetatable({}, {__mode = 'kv'})
local descriptions = setmetatable({}, {__mode = 'kv'})
local parameters = setmetatable({}, {__mode = 'k'})

local doc = {}

local decconcat = function(fn1, val) return fn1.fn(val) end
local deccall = function(fn, ...) return fn.fn(...) end
local function decorator(fn)
    local obj = {fn = fn}
    local mt = {
        __call = deccall,
        __concat = decconcat,
    }
    return setmetatable(obj, mt)
end

local function split(s, sep, n)
    if s == nil then
        error('missing string')
    end
    if sep == nil then
        error('missing separator')
    end

    local result = {}
    local i, j = string.find(s, sep)
    local count = 0
    while i > 0 and (not n or count+1 < n) do
        count = count + 1
        result[#result+1] = string.sub(s, 1, i-1)
        s = string.sub(s, j)
        i, j = string.find(s, sep)
    end
    result[#result+1] = s
    return result
end

local function load_docs(val)
    local sig = signatures[val]
    local desc = descriptions[val]
    local params = parameters[val]
    if sig == nil and desc == nil and params == nil then
        return nil
    end
    return {
        sig = sig,
        desc = desc,
        params = params,
    }
end

local _sig = function(sig)
    return decorator(function(fn)
        signatures[fn] = sig
        return fn
    end)
end
local _desc =  function(desc)
    return decorator(function(fn)
        descriptions[fn] = desc
        return fn
    end)
end
local _param = function(param)
    local pieces = split(param, '%s+', 2)
    local name = pieces[1]
    local desc = pieces[2] or ''

    return decorator(function(fn)
        local p = parameters[fn]
        if p == nil then p = {} end
        p[name] = desc
        parameters[fn] = p
        return fn
    end)
end

doc.sig =
    _sig[[s => fn => fn]] ..
    _desc[[A decorator that documents a function's signature.]] ..
    _param[[s   String containing the function signature]] ..
    _param[[fn  Function being documented]] ..
    _sig

doc.desc =
    _sig[[d => fn => fn]] ..
    _desc[[A decorator that documents a function's description.]] ..
    _param[[d   String containing the function description]] ..
    _param[[fn  Function being documented]] ..
    _desc

doc.param =
    _sig[[p => fn => fn]] ..
    _desc[[A decorator that documents a function parameter.]] ..
    _param[[p   String with name and description separated by whitespace]] ..
    _param[[fn  Function being documented]] ..
    _param

doc.help =
    doc.sig[[val =>  ()]] ..
    doc.desc[[Display help for an object, writing it to standard output]] ..
    doc.param[[val  Any table or function]] ..
    function(val)
        local d = load_docs(val)
        if d == nil then
            return
        end
        if d.desc then
            print(d.desc)
        end
        if d.sig then
            print()
            print(string.format('  %s', d.sig))
        end
        if d.params then
            print()
            for name, desc in pairs(d.params) do
                print(string.format('  %s\t%s', name, desc))
            end
        end
    end

    return doc
`
