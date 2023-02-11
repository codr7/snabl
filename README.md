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

`trace` may be used to toggle tracing of VM ops.

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
`test expected exp` evaluates `exp` and compares the result with `expected`.

```
  load "test/all.sl"

Testing T...OK
Testing F...OK
Testing 3...OK
Testing 1...OK
...

## Benchmarking
`bench n expr` evaluates `exp` `n` times and returns elapsed time.

```
  load "bench/fib.sl"

4.879072223s
```