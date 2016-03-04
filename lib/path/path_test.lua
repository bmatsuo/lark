local path = require('path')

function test_base()
    assert(path.base('/abc/def') == 'def')
    assert(path.base('/abc/') == 'abc')
    assert(path.base('/abc') == 'abc')
    assert(path.base('/') == '/')
end

function test_dir()
    assert(path.dir('/abc/def') == '/abc')
    assert(path.dir('/abc/') == '/abc')
    assert(path.dir('/abc') == '/')
    assert(path.dir('/') == '/')
end

function test_exists()
    assert(path.exists('path_test.lua'))
    assert(path.exists('path_test.go'))
    assert(path.exists('.'))
    assert(path.exists('..'))
    assert(not path.exists('./x/y/z'))
end

function test_ext()
    assert(path.ext('/x/y/z.zip') == '.zip')
    assert(path.ext('abc.tar.gz') == '.gz')
    assert(path.ext('/abc') == '')
    assert(path.ext('/abc.') == '.')
    assert(path.ext('/') == '')
end

function test_glob()
    local files = path.glob('./*.go')
    assert(table.getn(files) == 2)
    files = path.glob('../*/path.go')
    assert(table.getn(files) == 1)
    files = path.glob('x/y/*.z')
    assert(table.getn(files) == 0)
end

function test_is_dir()
    assert(not path.is_dir('./x/y/z'))
    assert(not path.is_dir('path_test.lua'))
    assert(not path.is_dir('path_test.go'))
    assert(path.is_dir('.'))
    assert(path.is_dir('..'))
end

function test_join()
    assert(path.join('abc', 'def') == 'abc/def')
    assert(path.join('', 'def') == 'def')
    assert(path.join('abc', '', 'def') == 'abc/def')
    assert(path.join('abc', '/def') == 'abc/def')
end
