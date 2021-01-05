package interceptors

import (
	"luis/app/controllers"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

func CheckAccess(c *revel.Controller) revel.Result {
	tok, err := c.Session.Get("access-token")
	if err != nil {
		// TODO: flash: you need to login?
		return c.Redirect(controllers.Access.LoginGet)
	}

	strTok, ok := tok.(string)
	if !ok {
		if err := c.Session.Set("access-token", ""); err != nil {
			revel.AppLog.Errorf("could clear session token: %q", err.Error())
		}
	}

	// Check if token is in cache.
	var secretTok struct{}

	if err := cache.Get(strTok, &secretTok); err != nil {
		// Token is not in cache or other error -> redirect to gateway.
		// TODO: flash: invalid credentials
		return c.Redirect(controllers.Access.LoginGet)
	}

	// Token is valid.

	return nil
}
