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

function test_named()
	local called = false
	local t =
		task.name[[task1]] ..
		function() called = true end

	assert(not task.find('t'))
	assert(task.find('task1'))
	task.find('task1')()
	assert(called)
end
