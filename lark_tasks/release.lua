local task = require('lark.task')
local path = require('path')
local version = require('version')
local moses = require('moses')

release = task .. function()
    lark.run('gen')
    lark.run('test')

    local release_root = 'release'
    local vx = version.get()
    vx = string.gsub(vx, '%W', '_')
    local name = 'lark-' .. vx
    local dist_template = name .. '-{{.OS}}-{{.Arch}}'
    local release_dir = path.join(release_root, name)
    local path_template = path.join(release_dir, dist_template, '{{.Dir}}')
    lark.exec{'mkdir', '-p', 'release'}
    lark.exec{'gox', '-os=!plan9', '-output='..path_template, '-ldflags='..ldflags, './cmd/...'}
    local dist_pattern = path.join(release_dir, '*')
    local dist_dirs = path.glob(dist_pattern)
    local ext_is = function(fp, ext) return path.ext(fp) == ext end
    dist_dirs = moses.reject(dist_dirs, function(_, dist) return ext_is(dist, '.zip') end)
    dist_dirs = moses.reject(dist_dirs, function(_, dist) return ext_is(dist, '.gz') end)
    for i, dist in pairs(dist_dirs) do
        lark.exec{'cp', 'README.md', 'CHANGES.md', 'LICENSE', 'AUTHORS', dist}
        lark.exec{'cp', '-r', 'docs', dist}

        local name = path.base(dist)
        if string.find(name, 'darwin') or string.find(name, 'windows') then
            local files = path.glob(path.join(dist, '*'))
            for i, fp in pairs(files) do
                files[i] = string.sub(fp, string.len(release_dir)+2)
            end
            lark.exec{'zip', '-o', name .. '.zip', files,
                      dir=release_dir}
        else
            lark.exec{'tar', '-cvzf', name .. '.tar.gz', name,
                      dir=release_dir}
        end
        lark.exec{'rm', '-r', dist}
    end
end
