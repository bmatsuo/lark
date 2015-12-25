--require 'lark'

lark.task{'demo', function ()
    lark.exec{'echo', 'an independent task', group=0} 

    lark.exec{'echo', group='setup', 'X'}

    build = lark.group{'build', after='setup'}
    lark.exec{'cc', '--version', group=build}

    lark.join{build}

    local cmd = {'wget', '-O', 'foo.txt', 'http://example.com/foo.txt'}
    lark.exec{cmd, group=0, ignore=true}

    lark.join()
end}

lark.task{'fail', function ()
    lark.exec{'false'}
end}
