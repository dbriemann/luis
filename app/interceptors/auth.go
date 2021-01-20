package interceptors

import (
	"luis/app/controllers"
	"luis/app/globals"
	"luis/app/models"
	"luis/app/store"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

func CheckAccess(c *revel.Controller) revel.Result {
	tok, err := c.Session.Get("access-token")
	if err != nil {
		c.Log.Debugf("could not get session token: %q", err.Error())
		c.Flash.Error(globals.ErrNotLoggedIn.Error())

		return c.Redirect(controllers.Access.Login)
	}

	strTok, ok := tok.(string)
	if !ok {
		c.Log.Debugf("could extract email for session token: %v", tok)
		if err := c.Session.Set("access-token", ""); err != nil {
			revel.AppLog.Errorf("could clear session token: %q", err.Error())
		}
	}

	// Check if token is in cache, and fetch corresponding mail.
	var email string

	if err := cache.Get(strTok, &email); err != nil {
		// Token is not in cache or other error -> redirect to gateway.
		c.Flash.Error(globals.ErrNotLoggedIn.Error())
		c.Log.Debugf("did not find token in cache: %s", strTok)

		return c.Redirect(controllers.Access.Login)
	}

	// Token is valid. Fetch user from DB and save in controller Args.
	user, err := models.UserByEmail(store.DB, email)
	if err != nil {
		c.Flash.Error(globals.ErrInternalServerError.Error())
		c.Log.Errorf("unexpected error: %q", err.Error())

		return c.RenderError(globals.ErrInternalServerError)
	}

	c.Args["user"] = *user

	adminEmail, _ := revel.Config.String("admin.email")
	if email == adminEmail {
		c.ViewArgs["IsAdmin"] = true
	} else {
		c.ViewArgs["IsAdmin"] = false
	}

	return nil
}
