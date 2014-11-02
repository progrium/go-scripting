package main

import (
	"fmt"
	"log"

	"github.com/progrium/go-extensions"
	"github.com/progrium/go-scripting"
	"github.com/progrium/go-scripting/ottojs"
)

var observers = extensions.ExtensionPoint(new(ProgramObserver))

type ProgramObserver struct {
	ProgramStarted func()
	ProgramFinished func()
}

func main() {
	ottojs.Register()
	scripting.UpdateGlobals(map[string]interface{}{
		"println": fmt.Println,
		"implements": func(module, iface string) {
			proxy := extensions.NewProxy(iface, 
				func (method string, args []interface{}) interface{} {
					value, err := scripting.Call(module, method, args)
					if err != nil {
						log.Println("error calling into", module, "with", method)
						return nil
					}
					return value
				},
			)
			extensions.RegisterWithName(iface, proxy, module)
		},
	})
	scripting.LoadModulesFromPath(".")

	for _, observer := range observers.All() {
		observer.(ProgramObserver).ProgramStarted()
	}
	
	fmt.Println("NORMALLY A PROGRAM DOES STUFF HERE")

	for _, observer := range observers.All() {
		observer.(ProgramObserver).ProgramFinished()
	}
}