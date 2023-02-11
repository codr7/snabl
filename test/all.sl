test T = 1 1
test F = 1 2

test T > 2 1
test F > 1 2

test 3 + 1 2
test 1 - 3 2

defun fib(n)
  if > n 1 + fib - n 1 fib - n 2 else n

test 55 fib 10

defun fibt(n a b)
  if > n 1 fibt - n 1 b + a b else if = n 0 a else b

test 55 fibt 10 0 1
