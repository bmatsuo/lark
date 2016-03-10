local lark = require('lark')

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
