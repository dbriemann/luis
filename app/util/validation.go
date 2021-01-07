package util

import "regexp"

// W3C email validation regex.
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsEMailValid returns true if the mail is valid after W3C and false otherwise.
func IsEMailValid(m string) bool {
	if len(m) < 3 || len(m) > 254 {
		return false
	}

	return emailRegex.MatchString(m)
}
