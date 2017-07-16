package commandMatchers

import (
	"regexp"
	"strings"
)

// MatchQuitCommand - Match quit command
func MatchQuitCommand(command string) bool {
	trimmedCommand := strings.Trim(command, " ")
	switch trimmedCommand {
	case "q", "quit", ":q":
		return true
	default:
		return false
	}
}

// MatchSearchCommand - Match search command
func MatchSearchCommand(command string) bool {
	trimmedCommand := strings.Trim(command, " ")
	matched, err := regexp.MatchString("^(s|search) .+", trimmedCommand)
	if err != nil {
		panic(err)
	}
	return matched
}

// MatchDownloadCommand - Match download command
func MatchDownloadCommand(command string) bool {
	trimmedCommand := strings.Trim(command, " ")
	matched, err := regexp.MatchString("^d [0-9]+$", trimmedCommand)
	if err != nil {
		panic(err)
	}
	return matched
}

// MatchOpenFolderCommand - Match search command
func MatchOpenFolderCommand(command string) bool {
	trimmedCommand := strings.Trim(command, " ")
	matched, err := regexp.MatchString("^(o|open)$", trimmedCommand)
	if err != nil {
		panic(err)
	}
	return matched
}
