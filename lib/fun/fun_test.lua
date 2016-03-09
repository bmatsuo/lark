fun = require('fun')

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
end

function test_map()
	local a = nil
	a = fun.map({2, 4, 6}, iseven)
	assert(equal(a, {false, true, false}))
	a = fun.map({2, 4, 6}, isodd)
	assert(equal(a, {true, false, true}))
end

function test_vmap()
	local a = nil
	a = fun.vmap({2, 4, 6}, iseven)
	assert(equal(a, {true, true, true}))
	a = fun.vmap({2, 4, 6}, isodd)
	assert(equal(a, {false, false, false}))
end

function test_select()
	local a = nil
	a = fun.sel({2, 4, 6}, iseven)
	assert(equal(a, {4}))
	a = fun.sel({2, 4, 6}, isodd)
	assert(equal(a, {2, 6}))
end

function test_vselect()
	local a = nil
	fun.vsel({1,2,3}, function() return true end)
	a = fun.vsel({2, 4, 6}, iseven)
	assert(equal(a, {2, 4, 6}))
	a = fun.vsel({2, 4, 6}, isodd)
	assert(equal(a, {}))
end
