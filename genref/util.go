package main

import "fmt"

const (
	CRED    = "\033[31m"
	CGREEN  = "\033[32m"
	CYELLOW = "\033[33m"
	CEND    = "\033[0m"
)

// perror prints message only when verbose mode is turned on
func perror(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(CRED+format+CEND+"\n", a...)
}

// pwarning prints message only when verbose mode is turned on
func pwarning(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(CYELLOW+format+CEND+"\n", a...)
}

// psuccess prints message only when verbose mode is turned on
func psuccess(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(CGREEN+format+CEND+"\n", a...)
}

// psuccess prints message only when verbose mode is turned on
func pinfo(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(format+"\n", a...)
}

// pverbose prints message only when verbose mode is turned on
func pverbose(format string, a ...interface{}) (n int, err error) {
	if *flVerbose {
		return fmt.Printf(format+"\n", a...)
	}
	return 0, nil
}

// containsString checks if a given string is a member of the string list
func containsString(sl []string, str string) bool {
	for _, s := range sl {
		if str == s {
			return true
		}
	}
	return false
}
