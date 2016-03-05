local doc = require('doc')

local version =
    doc.desc[[
    A collection utility functions for working with the lark project version.
    ]] ..
    doc.var[[filename  string
    -- Explicitly specify a default file path containing the version.  If no
    path is specified the default path "VERSION" is used.
    ]] ..
    {}

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

version.get =
    doc.desc[[
    Return the current project version by reading the version file.  The return
    value is cached so future calls to get() do not need to re-read the file.
    ]] ..
    doc.param[[filename  string
    -- The file to read the version number from.  When not given the default
    filename configured for the module will be read.
    ]] ..
    function(filename)
        filename = get_filename(filename)
        if not version_string or filename ~= version_file then
            version_file = filename
            version_string = version.read(filename)
        end
        return version_string
    end

version.read =
    doc.desc[[
    Return the current project version by reading the version file.  The read()
    function will read the version file every time it is called, unlike the
    get() function.
    ]] ..
    doc.param[[filename  string
    -- The file to read the version number from.  When not given the default
    filename configured for the module will be read.
    ]] ..
    function (filename)
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
