package globals

import (
	"fmt"
	"time"
)

const (
	CookieTokenBytes    = 128
	OneWeek             = 168 * time.Hour
	MaxSecretCollisions = 10
	ThumbnailSizePixels = 400
	GalleryColumns      = 4
)

var (
	ErrInvalidEmail        = fmt.Errorf("Please enter a valid email.")
	ErrInternalServerError = fmt.Errorf("Sorry this happened! The incident was logged.")
	ErrNotLoggedIn         = fmt.Errorf("You need to be logged in to access that page.")
	ErrStorageUnavailable  = fmt.Errorf("Storage is currently unavailable. Please contact the administrator.")
)
