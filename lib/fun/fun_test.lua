local fun = require('fun')
local math = require('math')

local function equal(a, b)
	if type(a) ~= type(b) then
		return false
	end
	if type(a) ~= 'table' then
		return a == b
	end
	if #a ~= #b then
		return false
	end
	for k, v in pairs(a) do
		if not equal(b[k], v) then
			return false
		end
	end
	for k, v in pairs(b) do
		if not equal(a[k], v) then
			return false
		end
	end
	return true
end

local function iseven(x) return x % 2 == 0 end
local function isodd(x) return x % 2 == 1 end

function test_equal()
	assert(equal(1, 1))
	assert(equal({}, {}))
	assert(equal({1, 2, {3}}, {1, 2, {3}}))
	assert(not equal(1, 2))
	assert(not equal({}, 2))
	assert(not equal({1, 2, {3}}, {1, 2, {3, 4}}))

	local a = {}
	a[1] = 1
	a[3] = 3
	assert(equal(a, {1, nil, 3}))
end

function test_flatten()
	local a = nil
	a = fun.flatten({1, 2, {3}})
	assert(equal(a, {1, 2, 3}))
	a = fun.flatten({1, 2, {3}}, 0)
	assert(equal(a, {1, 2, {3}}))
	a = fun.flatten({1, 2, {3, {4, 5}}}, 1)
	assert(equal(a, {1, 2, 3, {4, 5}}))
	a = fun.flatten({1, 2, {3, {4, 5}}}, -1)
	assert(equal(a, {1, 2, 3, 4, 5}))
	a = fun.flatten({1, 2, {3, {4, 5}}})
	assert(equal(a, {1, 2, 3, 4, 5}))
	a = fun.flatten({1, 2, {3, array={4, 5}}})
	assert(equal(a, {1, 2, 3}))
end

function test_map()
	local a = nil
	a = fun.map({2, 4, 6}, iseven)
	assert(equal(a, {false, true, false}))
	a = fun.map({2, 4, 6}, isodd)
	assert(equal(a, {true, false, true}))
	a = {2, 4, 6}
	a = fun.map(a, function(i, v) return #a-i+1, v end)
	assert(equal(a, {6, 4, 2}))

	a = {}
	a[1] = 2
	a[3] = 6
	a = fun.map(a, function() return 1 end)
	assert(equal(a, {1, nil, 1}))

	-- a really convoluted max computation
	local max = nil
	a = fun.map({1, 2, 3, a=1, b=2, c=3}, function(k, v)
		if max then
			max = math.max(max, v)
		else
			max = v
		end
		return 'max', max
	end)
	assert(equal(a, {max = 3}))
end

function test_vmap()
	local a = nil
	a = fun.vmap({2, 4, 6}, iseven)
	assert(equal(a, {true, true, true}))
	a = fun.vmap({2, 4, 6}, isodd)
	assert(equal(a, {false, false, false}))

	a = {}
	a[1] = 2
	a[3] = 6
	a = fun.vmap(a, function() return 1 end)
	assert(equal(a, {1, nil, 1}))

	a = fun.vmap({1, 2, 3, a=1, b=2, c=3}, function() return 1 end)
	assert(equal(a, {1, 1, 1, a=1, b=1, c=1}))
end

function test_select()
	local a = nil
	a = fun.sel({2, 4, 6}, iseven)
	assert(equal(a, {4}))
	a = fun.sel({2, 4, 6}, isodd)
	assert(equal(a, {2, 6}))

	a = {}
	a[1] = 2
	a[3] = 6
	a = fun.sel(a, function() return true end)
	assert(equal(a, {2, 6}))

	a = fun.sel({1, 2, 3, a=1, b=2, c=3}, function(k) return type(k) == 'string' end)
	assert(equal(a, {a=1, b=2, c=3}))

	a = fun.sel({1, 2, 3, a=1, b=2, c=3}, function(k) return type(k) == 'number' end)
	assert(equal(a, {1, 2, 3}))
end

function test_vselect()
	local a = nil
	fun.vsel({1,2,3}, function() return true end)
	a = fun.vsel({2, 4, 6}, iseven)
	assert(equal(a, {2, 4, 6}))
	a = fun.vsel({2, 4, 6}, isodd)
	assert(equal(a, {}))

	a = {}
	a[1] = 2
	a[3] = 6
	a = fun.vsel(a, function() return true end)
	assert(equal(a, {2, 6}))

	a = fun.vsel({1, 2, 3, a=1, b=2, c=3}, isodd)
	assert(equal(a, {1, 3, a=1, c=3}))

	a = fun.vsel({1, 2, 3, a=1, b=2, c=3}, iseven)
	assert(equal(a, {2, b=2}))
end
