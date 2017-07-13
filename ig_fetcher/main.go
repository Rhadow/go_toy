package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ahmdrz/goinsta"
	"github.com/howeyc/gopass"
)

// SimplifiedUser - Simplified twitter user data structure
type SimplifiedUser struct {
	ID        int64
	Name      string
	Username  string
	Followers int
}
type byFollowersCount []SimplifiedUser

func (s byFollowersCount) Len() int {
	return len(s)
}
func (s byFollowersCount) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byFollowersCount) Less(i, j int) bool {
	return s[i].Followers > s[j].Followers
}

var simplifiedUsers byFollowersCount

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
	fmt.Printf(loadingText)
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

func matchDownloadCommand(command string) bool {
	trimmedCommand := strings.Trim(command, " ")
	matched, err := regexp.MatchString("^d [0-9]+", trimmedCommand)
	if err != nil {
		panic(err)
	}
	return matched
}

func searchUserHandler(query string, insta *goinsta.Instagram) {
	users, searchErr := insta.SearchUsername(query)
	if searchErr != nil {
		panic(searchErr)
	}
	simplifiedUsers = byFollowersCount{}
	for _, user := range users.Users {
		simplifiedUsers = append(simplifiedUsers, SimplifiedUser{
			Name:      user.FullName,
			Username:  user.Username,
			Followers: user.FollowerCount,
			ID:        user.Pk,
		})
	}
	for i, user := range simplifiedUsers {
		fmt.Printf("[%d] %s - %s - %d followers\n", i+1, user.Username, user.Name, user.Followers)
	}
}

func downloadPictureHandler(targetIndex int, insta *goinsta.Instagram) {
	realIndex := targetIndex - 1
	if len(simplifiedUsers) <= realIndex || realIndex < 0 {
		fmt.Printf("%d out of range\n", targetIndex)
		return
	}
	targetUser := simplifiedUsers[realIndex]
	feeds, _ := insta.UserFeed(targetUser.ID, "", "")
	// feeds.Items[0].ImageVersions2.Candidates[0].URL
	fmt.Print(feeds)
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
			trimmedCommand := strings.Trim(command, " ")
			commandAndArgs := strings.Split(trimmedCommand, " ")
			query := strings.Join(commandAndArgs[1:], " ")
			searchUserHandler(query, insta)
		case matchDownloadCommand(command):
			trimmedCommand := strings.Trim(command, " ")
			commandAndArgs := strings.Split(trimmedCommand, " ")
			targetIndex, atoiErr := strconv.Atoi(commandAndArgs[1])
			if atoiErr != nil {
				panic(atoiErr)
			}
			downloadPictureHandler(targetIndex, insta)
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
	numberOfWorkers := flag.Int("w", 2, "Number of workers")
	flag.Parse()

	username, password := promptUserCredentials()
	insta := loginInstagram(username, password)
	defer insta.Logout()

	startTerminalInteraction(insta)
	quittingPrompt()
	fmt.Printf("Number of workers: %d\n", *numberOfWorkers)
}
