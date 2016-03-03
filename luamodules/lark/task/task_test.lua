local task = require('lark.task')

function test_create()
	local called = 0
	fn1 = function() called = called + 1 end
	anon_task1 =
		task.create ..
		function() fn1() end
	assert(not task.find('fn1'))
	assert(task.find('anon_task1'))
	task.find('anon_task1')()
	assert(called)
end

function test_module()
	local called = 0
	fn2 = function() called = called + 1 end
	anon_task2 =
		task ..
		function() fn2() end
	assert(not task.find('fn2'))
	assert(task.find('anon_task2'))
	task.find('anon_task2')()
	assert(called)
end

function test_with_name()
	local called = false
	local t =
		task.with_name[[task1]] ..
		function() called = true end

	assert(not task.find('t'))
	assert(task.find('task1'))
	task.find('task1')()
	assert(called)
end

function test_with_pattern()
	local called_svg = false
	local called_png = false
	task.with_pattern[[.*%.svg$]](function() called_svg = true end)
	task.with_pattern[[.*%.png$]](function() called_png = true end)

	assert(not task.find('foo.txt'))
	assert(task.find('foo.png'))
	assert(task.find('foo.svg'))
	task.find('foo.svg')()
	assert(called_svg)
	assert(not called_png)
end

function test_get_name()
	assert(task.get_name() == nil)
	assert(task.get_name({}) == nil)
	assert(task.get_name({name = 'x'}) == 'x')
	assert(task.get_name({pattern = 'x'}) == nil)
end

function test_get_pattern()
	assert(task.get_pattern() == nil)
	assert(task.get_pattern({}) == nil)
	assert(task.get_pattern({name = 'x'}) == nil)
	assert(task.get_pattern({pattern = 'x'}) == 'x')
end

function test_get_param()
	assert(task.get_param({}, "abc") == nil)
	assert(task.get_param({name = 'x'}, "abc") == nil)
	assert(task.get_param({param = 'x'}, "abc") == nil)
	assert(task.get_param({params = {abc = 'x'}}, "abc") == 'x')
	assert(task.get_param({params = {def = 'x'}}, "abc") == nil)
end

function test_run()
	local tpatt = 'run_test_pattern_*'
	local gotpatt = nil
	local t = function(ctx) gotpatt = task.get_pattern(ctx) end
	task.with_pattern(tpatt)(t)
	task.run('run_test_pattern_foo')
	assert(gotpatt)
	assert(gotpatt == tpatt)
end
