# POSIX-like Shell in Go

A lightweight command-line shell implementation written in Go.

## Features

* **REPL Loop:** Interactive prompt (`$ `) with standard input parsing and line trimming.
* **Builtin Commands:**
  * `exit`: Gracefully exits the shell.
  * `echo`: Prints arguments back to standard output.
  * `type`: Checks whether a given command is a shell builtin or resolves its full binary path using `$PATH`.
  * `pwd`: Prints the current working directory.
  * `cd`: Changes the current working directory.
* **External Command Execution:** Locates and executes external binaries found in `$PATH` using `os/exec`, streaming `stdout` and `stderr` back to the terminal.
* **PATH Resolution:** Manual PATH inspection to inspect executable permissions (`0111`) across directory splits.

## Getting Started

### Prerequisites

* Go 1.20 or later installed on your system.

### Running locally

1. Clone the repository:
```bash
   git clone https://github.com/jim3/m4z3y-5h3ll.git
   cd m4z3y-5h3ll
```

2. Run the shell directly:
```bash
go run main.go
```


3. Or compile and execute the binary:
```bash
go build -o m4z3y-5h3ll main.go
./m4z3y-5h3ll
```

## Usage Example

```bash
$ ./m4z3y-5h3ll
$ echo Hello, World!
Hello, World!
$ type echo
echo is a shell builtin
$ type ls
ls is /usr/bin/ls
$ type nonexistent
nonexistent: not found
$ exit
```
