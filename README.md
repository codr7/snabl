# Snabl
Snabl is a simple Go scripting language.

## Setup

```
git clone https://github.com/codr7/snabl.git
cd snabl
go build main/snabl.go
```
```
./snabl
  say "hello"

hello
```
```
./snabl help
Snabl v4

Usage:
snabl [command] [file1.sl] [file2.sl]...

Commands:
eval	Evaluate code and exit
read	Dump forms and exit
emit	Dump code and exit
repl	Evaluate code and start REPL
```

## Syntax
Snabl uses strict prefix syntax with optional parens; as a consequence, functions and macros have fixed arity.

## Types

### Bool
The type of boolean values.

```
  = T F

F
```

### Form
The type of source code forms.

### Fun
The type of functions.

```
  defun foo(x) x
  foo 42

42
```

### Int
The type of integer values.

### Macro
The type of macros.

### Meta
The type of types.

### Pos
The type of source code positions.

```
  1 2 3 pos
  
repl@1:7
```

### Prim
The type of primitives, functions implemented in Go.

### Slice
The type of slices.

```
  len [1 2 3]

3
```

### String
The type of string values.

```
  len "foo"

3
```

### Time
The type of time intervals.

## Environments
Curlies may be used to create new compile time environments.

```
   {defun foo() 42 foo}

[42]
  foo

repl@1:1 Error: foo?
```

## Debugging
`debug` may be used to toggle generation of debug info and panic on errors.

```
  fail "failing"
  
Error: failing
  debug
  
  fail "failing"
  
panic: repl@1:6 Error: failing
```

`trace` may be used to toggle tracing of operations and stack contents.

```
  trace
  
1 STOP []
  1 2 3
  
3 PUSH_INT 1 []
5 PUSH_INT 2 [1]
7 PUSH_INT 3 [1 2]
9 STOP [1 2 3]
3
```

## Testing
`test expected expr` evaluates `expr` and compares the result with `expected`.

```
  test T = 1 1
  
T = 1 1
  test T = 1 2
  
T = 1 2
TEST FAILED repl@1:8 F
```

You may find more examples [here](https://github.com/codr7/snabl/blob/main/test/all.sl).

## Benchmarking
`bench n expr` evaluates `expr` `n` times and returns elapsed time.

```
  load "bench/fib.sl"

2.807621686s
```

```
  load "bench/fibt.sl"

1.358420848s
```