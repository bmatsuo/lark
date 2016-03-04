local decorator = require('decorator')

function test_create()
    local _wrap = function(x) return {val = x} end
    local wrap = decorator.create(_wrap)
    local x = 0

    local x = wrap(wrap(1))
    assert(x.val.val == 1)

    x = wrap .. wrap .. wrap .. wrap .. 2
    assert(x.val.val.val.val == 2)
end

function test_annotator()
    local tset = {}
    local set = decorator.annotator(tset)

    local x =
        set('a') ..
        set('b') ..
        set('c') .. 
        {}

    assert(tset[x])
    assert(tset[x] == 'a')

    local tprepend = {}
    local prepend = decorator.annotator(tprepend, true)

    x =
        prepend('a') ..
        prepend('b') ..
        prepend('c') .. 
        x

    assert(tprepend[x])
    assert(table.concat(tprepend[x]) == "abc")
end
