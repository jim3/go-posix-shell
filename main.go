package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
			exitCode := 0
			if len(args) > 1 {
				if n, err := strconv.Atoi(args[1]); err == nil {
					exitCode = n
				}
			}
			os.Exit(exitCode)
		case "echo":
			fmt.Println(strings.Join(args[1:], " "))

		case "type":
			if len(args) < 2 {
				continue
			}

			target := args[1]
			if _, ok := builtin[target]; ok {
				fmt.Println(target, "is a shell builtin")
				continue
			}

			if fullPath, ok := findExe(target); ok {
				fmt.Println(target, "is", fullPath)
				continue
			}
			fmt.Printf("%s: not found\n", target)
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
			if _, exists := builtin[cmd]; exists {
				continue
			}
			command := exec.Command(cmd, args[1:]...)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			if err := command.Run(); err != nil {
				if errors.Is(err, exec.ErrNotFound) {
					fmt.Printf("%s: command not found\n", cmd)
					continue
				}

				if _, lookPathErr := exec.LookPath(cmd); lookPathErr != nil {
					fmt.Printf("%s: command not found\n", cmd)
					continue
				}

				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func findExe(exe string) (string, bool) {
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
