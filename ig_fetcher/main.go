package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ahmdrz/goinsta"
	"github.com/howeyc/gopass"
)

const loadingText string = "Loading..."

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

func loginInstagram(username, password string) *goinsta.Instagram {
	insta := goinsta.New(username, password)
	fmt.Println(loadingText)
	if err := insta.Login(); err != nil {
		panic(err)
	}
	return insta
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

func matchSearchCommand(command string) bool {
	trimmedCommand := strings.Trim(command, " ")
	matched, err := regexp.MatchString("^(s|search) .+", trimmedCommand)
	if err != nil {
		panic(err)
	}
	return matched
}

func searchUserHandler(command string) {
	trimmedCommand := strings.Trim(command, " ")
	commandAndArgs := strings.Split(trimmedCommand, " ")
	targetUser := strings.Join(commandAndArgs[1:], " ")
	fmt.Println(targetUser)
}

func startTerminalInteraction(insta *goinsta.Instagram) {
	const prompt = "ig_fetcher> "
	res, err := insta.GetProfileData()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello! %s\n", res.User.FullName)
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	quit := false
	for scanner.Scan() {
		command := scanner.Text()
		switch {
		case matchQuitCommand(command):
			quit = true
		case matchSearchCommand(command):
			searchUserHandler(command)
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

func quittingPrompt(insta *goinsta.Instagram) {
	fmt.Println("Good bye!")
}

func main() {
	numberOfWorkers := flag.Int("w", 2, "Number of workers")
	flag.Parse()

	username, password := promptUserCredentials()
	insta := loginInstagram(username, password)
	defer insta.Logout()

	startTerminalInteraction(insta)
	quittingPrompt(insta)
	fmt.Printf("Number of workers: %d\n", *numberOfWorkers)
}
