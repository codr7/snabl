package main

import (
	"github.com/codr7/snabl"
)

func main() {
	var vm snabl.Vm
	vm.Init()

	snabl.Repl(&vm)
}
