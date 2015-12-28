package main

import (
	"log"

	"github.com/codegangsta/cli"
)

// CommandList implements the "list" action and prints available tasks to
// standard output.
var CommandList = Command(func(lark *Context, cmd *cli.Command) {
	cmd.Name = "list"
	cmd.Usage = "List lark project task(s)"
	cmd.Action = lark.Action(List)
})

// List loads a lua vm and prints all defined tasks to standard output.
func List(c *Context) {
	luaFiles, err := FindTaskFiles("")
	if err != nil {
		log.Fatal(err)
	}

	luaConfig := &LuaConfig{}
	c.Lua, err = LoadVM(luaConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Lua.Close()

	err = InitLark(c, luaFiles)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Lua.DoString(`
	local names = {}
	for name in pairs(lark.tasks) do table.insert(names, name) end
	table.sort(names)
	for _, name in pairs(names) do
		if lark.default_task == name then
			print('  ' .. name .. ' (default)')
		else
			print('  ' .. name)
		end
	end
	`)
	if err != nil {
		log.Fatal(err)
	}
}
