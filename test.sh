#!/bin/bash


assert() {
  expected="$1"
  input="$2"

  ./went "$input" > tmp.s
  cc -o tmp tmp.s
  ./tmp
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

assert 0 'main() { return 0; }'
assert 42 'main() { return 42; }'
assert 21 'main() { return 5 + 20 -4; }'
assert 41 'main() { return 12 + 34 - 5; }'
assert 4 'main() { return (3 + 5) / 2; }'
assert 10 'main() { return -10 + 20; }'
assert 10 'main() { return - - + 10; }'

assert 0 'main() { return 0 == 1; }'
assert 1 'main() { return 42 == 42; }'
assert 1 'main() { return 0 != 1; }'
assert 0 'main() { return 42 != 42; }'

assert 1 'main() { return 0 < 1; }'
assert 0 'main() { return 1 < 1; }'
assert 0 'main() { return 2 < 1; }'
assert 1 'main() { return 0 <= 1; }'
assert 1 'main() { return 1 <= 1; }'
assert 0 'main() { return 2 <= 1; }'

assert 1 'main() { return 1 > 0; }'
assert 0 'main() { return 1 > 1; }'
assert 0 'main() { return 1 > 2; }'
assert 1 'main() { return 1 >= 0; }'
assert 1 'main() { return 1 >= 1; }'
assert 0 'main() { return 1 >= 2; }'

assert 6 'main() { foo = 1; bar = 2 + 3; return foo + bar; }'
assert 1 'main() { return 1; return 2; }'

assert 3 'main() { if (0) return 2; return 3; }'
assert 2 'main() { if (3 > 2) return 2; else return 3; }'
assert 3 'main() { if (3 < 2) return 2; else return 3; }'
assert 4 'main() { if (3 < 2) return 2; else if (2 < 1) return 3; else return 4; }'
assert 4 'main() { if (3 < 2) return 2; else if (2 < 1) return 3; else return 4; }'

assert 55 'main() { j = 0; for (i = 0; i <= 10; i = i + 1) j = i + j; return j; }'
assert 3 'main() { for (;;) return 3; return 5; }'

assert 3 'main() { {1; {2;} return 3;} }'
assert 55 'main() { j = 0; for (i = 0; i <= 10; i = i + 1) {tmp = i + j; j = tmp;} return j; }'

assert 32 'main() { return ret32(); } ret32() { return 32; }'

assert 7 'main() { return add(3, 4); } add(x, y) { return x + y; }'
assert 1 'main() { return sub(4, 3); } sub(x, y) { return x - y; }'
assert 55 'main() { return fib(9); } fib(x) { if (x <= 1) { return 1; } return fib(x - 1) + fib(x - 2); }'
assert 3 'main() { x = 3; y = &x; return *y; }'

echo OK
