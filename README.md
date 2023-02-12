# Snabl

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

4.203362269s
```

```
  load "bench/fibt.sl"

2.141586283s
```