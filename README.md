# Snabl

## Tests

`test expected exp` evaluates `exp` and compares the result with `expected`.

```
  load "test/all.sl"

Testing T...OK
Testing F...OK
Testing 3...OK
Testing 1...OK
...

## Benchmarks

`bench n expr` evaluates `exp` `n` times and returns elapsed time.

```
  load "bench/fib.sl"
  
6.51501242s
```