# An interpreter written in Go

### Some notes on the features of the language:
- *Lexer*. Only single byte characters.
- *Parser*. Recursive decent.
    - Associates:
        - *Left*. Equality (`==, !=`), comparison (`>, >=, <, <=`), term (`+, -`), factor (`*, /`).
        - *Right*. Unary (`!, -`).

### Examples.

---
This is essentially a Go port of the Lox language interpreter
found in *Nystrom, B. (2015), Crafting Interpreters*.
