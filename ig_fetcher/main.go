package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

func promptUserCredentials() (string, string) {
	var (
		username      string
		password      string
		passwordError error
	)

	scanner := bufio.NewScanner(os.Stdin)
	for username == "" || username == "\n" {
		fmt.Printf("Username: ")
		scanner.Scan()
		username = scanner.Text()
	}

	for password == "" || password == "\n" {
		fmt.Printf("Password: ")
		var temp []byte
		temp, passwordError = gopass.GetPasswdMasked()
		if passwordError != nil {
			panic("Can't read password")
		}
		password = string(temp)
	}

	return username, password
}

func matchQuitCommand(command string) bool {
	trimmedCommand := strings.Trim(command, " ")
	switch trimmedCommand {
	case "q", "quit", ":q":
		return true
	default:
		return false
	}
}

func startTerminalInteraction() {
	const prompt = "ig_fetcher> "
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	quit := false
	for scanner.Scan() {
		command := scanner.Text()
		switch {
		case matchQuitCommand(command):
			quit = true
		default:
			if command != "" {
				fmt.Printf("Unknown command: %s\n", command)
			}
		}
		if quit {
			break
		}
		fmt.Print(prompt)
	}
}

func main() {
	numberOfWorkers := flag.Int("w", 2, "Number of workers")
	flag.Parse()

	username, password := promptUserCredentials()
	startTerminalInteraction()
	fmt.Printf("Number of workers: %d\n", *numberOfWorkers)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", password)
}
