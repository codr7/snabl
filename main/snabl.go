package main

import (
	"github.com/codr7/snabl"
)

func main() {
	var vm snabl.Vm
	vm.Init()
	vm.Env().Import(&vm.AbcLib, &vm, nil)
	snabl.Repl(&vm)
}
