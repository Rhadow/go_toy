package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ahmdrz/goinsta"
	"github.com/howeyc/gopass"
	"github.com/rhadow/go_toy/ig_fetcher/commandMatchers"
	"github.com/rhadow/go_toy/ig_fetcher/fetcherResponse"
	"github.com/rhadow/go_toy/ig_fetcher/instaHandlers"
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

func startTerminalInteraction(insta *goinsta.Instagram, simplifiedUsers *fetcherResponse.ByFollowersCount) {
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
		case commandMatchers.MatchQuitCommand(command):
			quit = true
		case commandMatchers.MatchSearchCommand(command):
			trimmedCommand := strings.Trim(command, " ")
			commandAndArgs := strings.Split(trimmedCommand, " ")
			query := strings.Join(commandAndArgs[1:], " ")
			instaHandlers.SearchUserHandler(query, insta, simplifiedUsers)
		case commandMatchers.MatchDownloadCommand(command):
			trimmedCommand := strings.Trim(command, " ")
			commandAndArgs := strings.Split(trimmedCommand, " ")
			targetIndex, atoiErr := strconv.Atoi(commandAndArgs[1])
			if atoiErr != nil {
				panic(atoiErr)
			}
			instaHandlers.DownloadPictureHandler(targetIndex, insta, simplifiedUsers)
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

func quittingPrompt() {
	fmt.Println("Good bye!")
}

func main() {
	var simplifiedUsers fetcherResponse.ByFollowersCount
	const loadingText string = "Loading..."

	numberOfWorkers := flag.Int("w", 2, "Number of workers")
	flag.Parse()

	username, password := promptUserCredentials()
	insta := instaHandlers.LoginInstagram(username, password, loadingText)
	defer insta.Logout()

	startTerminalInteraction(insta, &simplifiedUsers)
	quittingPrompt()
	fmt.Printf("Number of workers: %d\n", *numberOfWorkers)
}
