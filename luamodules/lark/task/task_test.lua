local task = require('lark.task')

function test_create()
	print("CREATE")
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
	print("MODULE")
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
	print("NAME")
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
	print("PATTERN")
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
