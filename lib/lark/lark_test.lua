local lark = require('lark')

function test_async()
    lark.start('true')
    assert(pcall(lark.wait))

    lark.start('false')
    assert(not pcall(lark.wait))

    lark.start('false', {ignore = true})
    assert(pcall(lark.wait))

    lark.start{'true'}
    assert(pcall(lark.wait))

    lark.start('false', {group = 'fail'})
    lark.start('true', {group = 'ok'})
    assert(pcall(lark.wait, 'ok'))
    assert(not pcall(lark.wait, 'fail'))
    assert(pcall(lark.wait, 'fail')) -- subsequent calls do not fail

    lark.start('false', {group = 'fail'})
    lark.start('true', {group = 'ok'})
    assert(not pcall(lark.wait))
    assert(pcall(lark.wait))
end

function test_exec()
    assert(pcall(lark.exec, {'true'}))
    assert(not pcall(lark.exec, {'false'}))

	local ok, out, err = false, nil, nil

    ok, out, err = pcall(lark.exec, {'true'})
	assert(ok)
	assert(not out)
	assert(not err)

    ok, out, err = pcall(lark.exec, {'false', ignore = true})
	assert(ok)
	assert(not out)
	assert(err)

    out, err = lark.exec('false', {ignore = true})
	assert(not out)
	assert(err)

	ok, out, err = pcall(lark. exec, {'echo', 'test output', stdout = '$'})
	assert(ok)
	assert(out)
	assert(out == 'test output\n')
	assert(not err)

	ok, out, err = pcall(lark. exec, 'echo', 'test output', {stdout = '$'})
	assert(ok)
	assert(out)
	assert(out == 'test output\n')
	assert(not err)
end
