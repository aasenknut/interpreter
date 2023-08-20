An interpreter written in Go

### Some notes on the features of the language:
- *Lexer*. Only single byte characters.
- *Parser*. Recursive decent.
    - Associates:
        - *Left*. Equality (`==, !=`), comparison (`>, >=, <, <=`), term (`+, -`), factor (`*, /`).
        - *Right*. Unary (`!, -`).

### Examples.

**Functions**.
```
fun sayHi(first, last) {
  print "Hi, " + first + " " + last + "!";
  return "Success!";
}

var x = sayHi("Dear", "Reader");
print x;

>> Hi, Dear Reader!
>> Success!
```

**Loops**.
```
var a = 0;
var temp;

for (var b = 1; a < 9; b = temp + b) {
  print a;
  temp = a;
  a = b;
}
```

---
This is essentially a Go port of the Lox language interpreter
found in *Nystrom, B. (2015), Crafting Interpreters*.
