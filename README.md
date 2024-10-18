# Nadia Lang
Nadia Lang is an event oriented programming language for Arudino inspired by Go, Python and JS.

# Install
Use *Go CLI* to install the Nadia CLI.
```
go install github.com/agustin-del-pino/nadia-lang/cmd/nadia
```

Once installed, you can use the Nadia CLI. Use the `nadia init` command to initilize the necessary files.
```
nadia init
```

# Nadia CLI

Use this command to display the help information.
```
nadia help

nadia <command> [param] ...flags


Commands

init:     initializes the internal dependencies. 
build:    builds the given nadia source file
help:     displays CLI's help
version:  displays nadia-lang's version


Flags

-o, --out:         sets the output filepath
-n, --nadia-path:  sets a temporally external lib path
```

# Hello World

Create a new `main.nad` file.

```
// include the native lib "serial"
include "nad:serial"

// add a listener function to the setup event
lst (setup) func on_init() {
    // set the bauds for serial communication
    set_bauds(9600)
    // print the message
    print("Hello World, from Nadia")
}
```

Then, build the code.

```
nadia build main.nad
```

# Features

| Feature                               | Represents                   | Status |
| ------------------------------------- | ---------------------------- | ------ |
| Variables                             | Var declaration              | ✅      |
| Constants                             | Const declaration            | ✅      |
| Functions                             | Function declaration         | ✅      |
| Objects                               | Struct declaration           | ✅      |
| Events                                | Event Init-Trigger Mechanism | ✅      |
| Listeners                             | Event listening Mechanism    | ✅      |
| Event Arguments Propagation           | Event-Args propagation       | ❌      |
| If                                    | Condition Statement          | ✅      |
| When                                  | Switch Statement             | ✅      |
| For                                   | While Statement              | ✅      |
| For-Range                             | For Statement                | ✅      |
| For-Each                              | For Each Statement           | ✅      |
| Return                                | Return Statement             | ✅      |
| No semicolon                          | Like Go, Python, etc         | ✅      |
| No parenthesis at statements          | Like Go, Python, etc         | ✅      |
| C/C++ Type Representation             | Represent type from C/C++    | ✅      |
| Type Check                            | Semantic Analyzer            | ❌      |
| Type Deduction                        | Like Go, TS, etc             | ❌      |
| Transpilation Macros                  | Similar to `#define`         | ✅      |
| Native Lib inclusion                  | Similar to `#include`        | ✅      |
| Custom Lib inclusion                  | Similar to `#include`        | ❌      |
| File inclusion                        | Similar to `#include`        | ✅      |
| File inclusion from other directories | Similar to `#include`        | ❌      |
| Syntax Highlight Plugins/Extension    | For IDE/Code Editors         | ❌      |
