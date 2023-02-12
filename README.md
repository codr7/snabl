# Snabl
Snabl is a simple Go scripting language.

## Setup

```
git clone https://github.com/codr7/snabl.git
cd snabl
go build main/snabl.go
./snabl
Snabl v4

  say "hello"

hello
```

## Syntax
Snabl uses strict prefix syntax with optional parens; as a consequence, functions and macros have fixed arity.

## Types

### Bool
The type of boolean values.

### Form
The type of source code forms.

### Fun
The type of functions.

### Int
The type of integer values.

### Macro
The type of macros.

### Meta
The type of types.

### Pos
The type of source code positions.

### Prim
The type of primitives, functions implemented in Go.

### String
The type of string values.

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
  
panic: repl@1:7 Error: failing
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

`pos` may be used to push the current source position.

```
  1 2 say pos 3
  
repl@1:9
3
```

## Testing
`test expected expr` evaluates `expr` and compares the result with `expected`.

```
  load "test/all.sl"

Testing T...OK
...
```

## Benchmarking
`bench n expr` evaluates `expr` `n` times and returns elapsed time.

```
  load "bench/fib.sl"

3.261481578s
```

```
  load "bench/fibt.sl"

1.821470141s
```