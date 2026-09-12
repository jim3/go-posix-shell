package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var builtin = map[string]bool{
	"echo":     true,
	"exit":     true,
	"type":     true,
	"cd":       true,
	"pwd":      true,
	"history":  true,
	"jobs":     true,
	"complete": true,
	"declare":  true,
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		line, err := reader.ReadString('\n')
		
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		args := strings.Fields(line)
		cmd := args[0]
		switch cmd {
		case "exit":
			os.Exit(0)

		case "echo":
			fmt.Println(strings.Join(args[1:], " "))

		case "type":
			if len(args) < 2 {
				fmt.Println("type: missing argument")
				continue
			}
			if _, ok := builtin[args[1]]; ok {
				fmt.Println(args[1], "is a shell builtin")
				continue
			}
			if fullPath, ok := execBin(args); ok {
				fmt.Println(args[1], "is", fullPath)
				continue
			} else {
				fmt.Printf("%s: not found\n", args[1])
			}
		case "pwd":
			pwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("%v: not found\n", pwd)
			}
			fmt.Println(pwd)
		case "cd":
			path := args[1]
			if path == "~" {
				path = os.Getenv("HOME")
			}
			err := os.Chdir(path)
			if err != nil {
				fmt.Println("cd: /non-existing-directory: No such file or directory")
			}
		default:
			_, exists := builtin[cmd]
			if !exists {
				_, err := exec.LookPath(cmd)
				if err != nil {
					fmt.Printf("%v: not found\n", cmd)
				} else {
					command := exec.Command(cmd, args[1:]...)
					command.Stdout = os.Stdout
					command.Stderr = os.Stderr
					output := command.Run()
					if output != nil {
						fmt.Println("Error:", output)
					}
				}
			}
		}
	}
}

func execBin(args []string) (string, bool) {
	if len(args) < 2 {
		return "", false
	}

	exe := args[1]
	for _, v := range filepath.SplitList(os.Getenv("PATH")) {
		filePath := filepath.Join(v, exe)
		info, err := os.Stat(filePath)
		
		if err == nil && !info.IsDir() && info.Mode().Perm()&0o111 != 0 {
			return filePath, true
		}

		if os.IsNotExist(err) {
			continue
		}
	}
	return "", false
}
