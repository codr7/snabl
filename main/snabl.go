package main

import (
	"fmt"
	"github.com/codr7/snabl"
	"log"
	"os"
)

func main() {
	var vm snabl.Vm
	vm.Init()
	vm.Env.Import(&vm.AbcLib, &vm, nil)

	args := os.Args[1:]
	var cmd string
	
	if len(args) > 0 {
		a := args[0] 
		if a == "eval" || a == "help" || a == "read" || a == "emit" || a == "repl" {
			cmd = a
			args = args[1:]
		} else {
			cmd = "eval"
		}
	} else {
		cmd = "repl"
	}

	if cmd == "help" {
		fmt.Printf("Snabl v%v\n\n", snabl.VERSION)
		fmt.Print("Usage:\n")
		fmt.Print("snabl [command] [file1.sl] [file2.sl]...\n\n")
		fmt.Print("Commands:\n")
		fmt.Print("eval\tEvaluate code and exit\n")
		fmt.Print("read\tDump forms and exit\n")
		fmt.Print("emit\tDump code and exit\n")
		fmt.Print("repl\tEvaluate code and start REPL\n")
	} else if cmd == "read" {
		var forms snabl.Forms

		for _, p := range args {
			if err := vm.LoadForms(p, &forms); err != nil {
				log.Fatal(err)
			}
		}

		for _, f := range forms.Items() {
			fmt.Println(f)
		}
	} else {
		for _, p := range args {
			if err := vm.Load(p, cmd == "eval" || cmd == "repl"); err != nil {
				log.Fatal(err)
			}
		}
		
		switch cmd {
		case "emit":
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
		default:
			log.Fatalf("%v?", cmd)
		}
	}
}
