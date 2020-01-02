// Golox implements the Lox programming language from
// https://craftinginterpreters.com.
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

var hadError = false

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: golox [script]")
		os.Exit(64) // EX_USAGE
	}
	if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(66) // EX_NOINPUT
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(66) // EX_NOINPUT
	}
	run(string(bytes))
}

func runPrompt() {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		run(scanner.Text())
		hadError = false
		fmt.Print("> ")
	}
	fmt.Println()
}

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()
	for i := 0; i < len(tokens); i++ {
		fmt.Println(tokens[i])
	}
}

func reportError(line int, message string, at string) {
	fmt.Fprintf(os.Stderr, "line %d: %s %s", line, message, at)
	hadError = true
}
