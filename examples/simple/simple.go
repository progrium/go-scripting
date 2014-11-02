package main

import (
	"fmt"

	"github.com/progrium/go-scripting"
	"github.com/progrium/go-scripting/ottojs"
)

func main() {
	ottojs.Register()
	scripting.LoadModulesFromPath(".")
	scripting.UpdateGlobals(map[string]interface{}{
		"println": fmt.Println,
	})
	scripting.Call("example", "helloworld", nil)
}