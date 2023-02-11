package snabl

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
)

func Repl(vm *Vm) {
	fmt.Printf("Snabl v%v\n\n", VERSION)	
	scanner := bufio.NewScanner(vm.Stdin)
	var buffer bytes.Buffer

	for {
	NEXT:
		fmt.Fprintf(vm.Stdout, "  ")
		
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			break
		}
		
		line := scanner.Text()

		if line == "" {
			pos := NewPos("repl", 1, 1)
			var forms Forms
			
			if err := ReadForms(vm, pos, bufio.NewReader(&buffer), &forms); err != nil {
				fmt.Fprintln(vm.Stdout, err)
				buffer.Reset()
				goto NEXT
			}

			buffer.Reset()
			pc := vm.EmitPc()

			for forms.Len() > 0 {
				if err := forms.Pop().Emit(&forms, vm, vm.Env()); err != nil {
					fmt.Fprintln(vm.Stdout, err)
					goto NEXT
				}
			}

			vm.Code[vm.Emit()] = StopOp()
			
			if err := vm.Eval(&pc); err != nil {
				fmt.Fprintln(vm.Stdout, err)
				vm.Stack.Clear()
				goto NEXT
			}
			
			if vm.Stack.Len() > 0 {
				fmt.Fprintln(vm.Stdout, *vm.Stack.Top(0))
				vm.Stack.Clear()
			}
		} else if _, err := fmt.Fprintln(&buffer, line); err != nil {
			log.Fatal(err)
		}
	}	
}
