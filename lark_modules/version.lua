local version = {}

local version_file = nil
local version_string = nil

local default_filename = 'VERSION'

version.filename = nil

local function get_filename(filename)
    if filename then
        return filename
    end
    if version.filename then
        return version.filename
    end
    return default_filename
end

version.get = function(filename)
    filename = get_filename(filename)
    if not version_string or filename ~= version_file then
        version_file = filename
        version_string = version.read(filename)
    end
    return version_string
end

version.read = function (filename)
    filename = get_filename(filename)

    local f, err = io.open(filename)
    if err then
        return nil, err
    end
    local v = f:read('*all')
    f:close()
    v = string.gsub(v, '#[^\n]\n', '')
    v = string.gsub(v, '%s', '')
    return v
end

return version
