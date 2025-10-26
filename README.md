# gosh
Go based shell which works cross platform (Linux, Windows), is interactive and intuitive to use.
- go is a compiled language which can compile instantly. we will leverage that to build an interactive shell which could hot reload. some existing solutions exist and we can brainstorm new ones as well if necessary.
- because of that we should be able to easily save some given commands as a script. in fact, it should by default write the shell "code" as a project in a specific folder.
- "exporting" the current session as a new script (cli) would be a new cobra cli inside cli/ and prompt user for a name
- shell commands should be aliased to go : typing them would write in our project the Go equivalent but every bash commands should still be translated and understood. Still all Go code should be valid.

- 
