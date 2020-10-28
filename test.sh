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

assert 0  0
assert 42 42
assert 21 "5+20-4"
assert 21 "5 + 20 - 4"
assert 26 "2 * 3 + 4 * 5"
assert 0  "4 / 2 - 10 / 5"
assert 20 "5 * ( 5 - 1 )"
assert 20 "+ 5 * + ( 5 - 1 )"
assert 20 "- 5 * - ( 5 - 1 )"

echo OK
