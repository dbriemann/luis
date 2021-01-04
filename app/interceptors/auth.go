package interceptors

import (
	"luis/app/controllers"

	"github.com/revel/revel"
)

func CheckAccess(c *revel.Controller) revel.Result {
	tok, err := c.Session.Get("access-token")
	if err != nil {
		return c.Redirect(controllers.Access.Login)
	}

	// TODO: check if access token is in DB.
	tokValid := false
	if tok != "" {
		// TODO
	}

	if !tokValid {
		return c.Redirect(controllers.Access.Login)
	}

	return nil
}
