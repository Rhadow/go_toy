package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/ahmdrz/goinsta"
	"github.com/howeyc/gopass"
	"github.com/rhadow/go_toy/ig_fetcher/commandMatchers"
	"github.com/rhadow/go_toy/ig_fetcher/fetcherResponse"
	"github.com/rhadow/go_toy/ig_fetcher/instaHandlers"
	"github.com/skratchdot/open-golang/open"
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

func startTerminalInteraction(insta *goinsta.Instagram, simplifiedUsers *fetcherResponse.ByFollowersCount, numberOfWorkers int, baseDir string) {
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
		trimmedCommand := strings.Trim(command, " ")
		commandAndArgs := strings.Split(trimmedCommand, " ")
		switch {
		case commandMatchers.MatchQuitCommand(trimmedCommand):
			quit = true
		case commandMatchers.MatchOpenFolderCommand(trimmedCommand):
			openErr := open.Run(baseDir)
			if openErr != nil {
				panic(openErr)
			}
		case commandMatchers.MatchSearchCommand(trimmedCommand):
			query := strings.Join(commandAndArgs[1:], " ")
			instaHandlers.SearchUserHandler(query, insta, simplifiedUsers)
		case commandMatchers.MatchDownloadCommand(trimmedCommand):
			targetIndex, atoiErr := strconv.Atoi(commandAndArgs[1])
			if atoiErr != nil {
				panic(atoiErr)
			}
			realIndex := targetIndex - 1
			if len(*simplifiedUsers) <= realIndex || realIndex < 0 {
				fmt.Printf("%d out of range\n", targetIndex)
				break
			}
			targetUser := (*simplifiedUsers)[realIndex]
			instaHandlers.DownloadPictureHandler(insta, targetUser, numberOfWorkers, baseDir)
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
	//Get system user folder
	usr, _ := user.Current()
	baseDir := fmt.Sprintf("%v/Pictures/goInstagram", usr.HomeDir)

	numberOfWorkers := flag.Int("w", 4, "Number of workers")
	flag.Parse()

	username, password := promptUserCredentials()
	insta := instaHandlers.LoginInstagram(username, password)
	defer insta.Logout()

	startTerminalInteraction(insta, &simplifiedUsers, *numberOfWorkers, baseDir)
	quittingPrompt()
}
