defun fibt(n a b)
  if > n 1 fibt - n 1 b + a b else if = n 0 a else b

say bench 100000 fibt 100 0 1