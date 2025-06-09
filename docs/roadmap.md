# Roadmap

1. Allow for custom marker definitions in the parser and lexer.
2. LocalLoader + HTTPLoader etc. (loader docs in general).
3. Make the loader much easier by just using the high level variables but more
   methods in the interface so the interaction is easier.
4. Variablen für marker (custom)
5. wildcards for typeid?
6. sdk package for save public usage
7. Build a variable in which the marker themselves will be included so
   referenceing other marker values is possible in the future e.g.
   path:to:i=32 -> $path:to:i resolves to 32.
