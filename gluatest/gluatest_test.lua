local module = require('test.module')

function helper()
    return 'HELP!'
end

function __test_setup()
	module.set_value('__test_setup', true)
end

function __test_teardown()
	module.set_value('__test_teardown', true)
end

function test_a()
	module.set_value('test_a', true)
end

function test_b()
	module.set_value('test_b', true)
end

function test_fail()
	module.set_value('test_fail', true)
	assert(false)
end

function test_ok()
	module.set_value('test_ok', true)
	assert(true)
end
