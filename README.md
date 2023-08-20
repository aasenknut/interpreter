An interpreter written in Go.

Some notes on the features of the language:
- *Lexer*. Only single byte characters.
- *Parser*. Recursive decent.
    - Associates:
        - *Left*. Equality (`==, !=`), comparison (`>, >=, <, <=`), term (`+, -`), factor (`*, /`).
        - *Right*. Unary (`!, -`).

---
This is essentially a port of the Lox languages interpreter, from Java to Go, 
found in *Nystrom, B. (2015), Crafting Interpreters*,
