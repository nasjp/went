#!/bin/bash

cat <<EOF | cc -xc -c -o tmp2.o -
int ret3() { return 3; }
int ret5() { return 5; }
int add(int x, int y) { return x+y; }
int sub(int x, int y) { return x-y; }
int add6(int a, int b, int c, int d, int e, int f) {
  return a+b+c+d+e+f;
}
EOF

assert() {
  expected="$1"
  input="$2"

  ./went "$input" > tmp.s
  cc -o tmp tmp.s tmp2.o
  ./tmp
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

assert 0  "0;"
assert 42 "42;"
assert 21 "5+20-4;"
assert 21 "5 + 20 - 4;"
assert 26 "2 * 3 + 4 * 5;"
assert 0  "4 / 2 - 10 / 5;"
assert 20 "5 * ( 5 - 1 );"
assert 20 "+ 5 * + ( 5 - 1 );"
assert 20 "- 5 * - ( 5 - 1 );"

assert 10 '- - 10;'
assert 10 '- - + 10;'

assert 0 '0 == 1;'
assert 1 '42 == 42;'
assert 1 '0 != 1;'
assert 0 '42 != 42;'

assert 1 '0 < 1;'
assert 0 '1 < 1;'
assert 0 '2 < 1;'
assert 1 '0 <= 1;'
assert 1 '1 <= 1;'
assert 0 '2 <= 1;'

assert 1 '1 > 0;'
assert 0 '1 > 1;'
assert 0 '1 > 2;'
assert 1 '1 >= 0;'
assert 1 '1 >= 1;'
assert 0 '1 >= 2;'

assert 1 'a = 1;'
assert 4 'a = 1; b = 3; c = a + b; c;'
assert 10 'aa = 1; bb = 3; aa = 7; aa + bb;'

assert 10 'a = 10; return a;'
assert 2 'return 2; return 3;'

assert 3 'if (0) return 2; return 3;'
assert 3 'if (3 < 2) return 2; else return 3;'
assert 4 'if (3 < 2) return 2; else if (2 < 1) return 3; else return 4;'

assert 55 'j = 0; for (i = 0; i <= 10; i = i + 1) j = i + j; return j;'
assert 3 'for (;;) return 3; return 5;'
assert 0 'for (i = 0;;)  if(i <= 10) return i;'

assert 3 '{ 1; {2;} return 3; }'

assert 55 '{ j = 0; for (i = 0; i <= 10; i = i + 1) {tmp = i + j; j = tmp;} return j; }'


assert 3 'return ret3();'
assert 5 'return ret5();'

assert 8 'return add(3, 5);'
assert 2 'return sub(5, 3);'
assert 21 'return add6(1,2,3,4,5,6);'

echo OK
