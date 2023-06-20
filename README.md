This is the *Glox* language.
(parsing expressions)

To say that this is inspired by *Nystrom, B. (2015), Crafting Interpreters*,
would be the understatement of the year.

Some notes on the features of the language:
- *Lexer*. Only single byte characters -- not UTF-8.
- *Parser*. Recursive decent.
    - Associates:
        - *Left*. Equality (`==, !=`), comparison (`>, >=, <, <=`), term (`+, -`), factor (`*, /`).
        - *Right*. Unary (`!, -`).

---
*Powered by technical debt*.
