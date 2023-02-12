package main

import (
	"github.com/codr7/snabl"
	"log"
	"os"
)

func main() {
	var vm snabl.Vm
	vm.Init()
	vm.Env().Import(&vm.AbcLib, &vm, nil)

	args := os.Args[1:]
	var cmd string
	
	if len(args) > 0 {
		if args[0] == "dump" || args[0] == "eval" || args[0] == "repl" {
			cmd = args[0]
			args = args[1:]
		} else {
			cmd = "eval"
		}
	} else {
		cmd = "repl"
	}

	for _, p := range args {
		if err := vm.Load(p, cmd != "dump"); err != nil {
			log.Fatal(err)
		}
	}
	
	switch cmd {
	case "dump":
		var pos *snabl.Pos
		
		for pc, op := range vm.Code {
			if op.Id() == snabl.POS_OP {
				p := vm.Tags[op.Pos()].Data().(snabl.Pos)
				pos = &p
			} else {
				op.Trace(&vm, pc, pos, false, vm.Stdout)
			}
		}
	case "eval":
	case "repl":
		snabl.Repl(&vm)
	}
}
