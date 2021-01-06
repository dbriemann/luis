package globals

import (
	"fmt"
	"time"
)

const (
	CookieTokenBytes    = 128
	OneWeek             = 168 * time.Hour
	MaxSecretCollisions = 10
)

var (
	ErrInternalServerError = fmt.Errorf("Sorry this happened! The incident was logged.")
	ErrNotLoggedIn         = fmt.Errorf("You need to be logged in to access that page.")
)
