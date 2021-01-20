package controllers

import (
	"crypto/subtle"
	"database/sql"
	"errors"

	"luis/app/globals"
	"luis/app/models"
	"luis/app/store"
	"luis/app/util"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

type Access struct {
	*revel.Controller
}

func (c Access) Login() revel.Result {
	return c.Render()
}

func (c Access) Logout() revel.Result {
	tok, err := c.Session.Get("access-token")
	if err == nil {
		strTok, ok := tok.(string)
		if ok {
			// We can safely ignore any deletion errors here.
			_ = cache.Delete(strTok)
		}
	}

	// If any failure happens before this point we can ignore it,
	// because then the user is not logged in anyways.

	c.Flash.Success("See you next time!")

	return c.Redirect(Access.Login)
}

func (c Access) LoginPost(remember bool) revel.Result {
	// TODO -- save email in flash/session? if redirecting to login again
	email := c.Params.Form.Get("email")
	pass := c.Params.Form.Get("password")

	c.Log.Infof("login attempt for email %q", email)

	user, err := models.UserByEmail(store.DB, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Log.Infof("login - unknown email: %q", email)
			c.Flash.Error("Invalid user credentials!")

			return c.Redirect(Access.Login)
		}
		// TODO - internal server error?
	}

	if subtle.ConstantTimeCompare([]byte(user.Secret), []byte(pass)) == 0 {
		c.Log.Infof("login - invalid password for email %q", email)
		c.Flash.Error("Invalid user credentials!")

		return c.Redirect(Access.Login)
	}

	var token string
	var tries int

	// User is now authed -> generate token.
	for {
		tries++

		// This should never collide with another token but better be sure..
		tok, err := util.GenerateSecret(globals.CookieTokenBytes)
		if err != nil {
			c.Log.Errorf("cannot generate secret: %q", err)

			return c.RenderError(globals.ErrInternalServerError)
		}

		// Save token in cache if it does not collide.
		if err := cache.Add(tok, email, globals.OneWeek); err != nil {
			if errors.Is(err, cache.ErrNotStored) {
				c.Log.Warnf("collision of secrets for email %q, try: %d", email, tries)

				if tries > globals.MaxSecretCollisions {
					// Something must be wrong. Just stop and fail.
					c.Log.Errorf("too many secret collisions")

					return c.RenderError(globals.ErrInternalServerError)
				}

				// Try another time.
				continue
			} else {
				c.Log.Errorf("unexpected error: %q", err.Error())

				return c.RenderError(globals.ErrInternalServerError)
			}
		}

		token = tok

		break
	}

	// Save token in session.
	if err := c.Session.Set("access-token", token); err != nil {
		c.Log.Errorf("could not write token to session: %q", err.Error())

		return c.RenderError(globals.ErrInternalServerError)
	}

	if remember {
		c.Session.SetDefaultExpiration()
	} else {
		c.Session.SetNoExpiration()
	}

	// Success -> redirect to home.
	return c.Redirect(App.Index)
}
