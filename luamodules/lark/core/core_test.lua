local core = require('lark.core')

function test_log()
    core.log{'test_log', color='green'}
end

function test_environ()
    local env = core.environ()
    assert(env.PATH)
end

function test_exec()
    local result = core.exec{'false'}
    assert(result.error)
    result = core.exec{'true'}
    assert(not result.error)
end
