package controllers

import (
	"crypto/subtle"
	"errors"
	"time"

	"luis/app/gormdb"
	"luis/app/models"
	"luis/app/util"

	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"gorm.io/gorm"
)

type Access struct {
	*revel.Controller
}

func (c Access) LoginGet() revel.Result {
	return c.Render()
}

func (c Access) LoginPost(remember bool) revel.Result {
	email := c.Params.Form.Get("email")
	pass := c.Params.Form.Get("password")

	c.Log.Infof("login attempt for email %q", email)

	var user models.User

	result := gormdb.DB.Take(&user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Log.Infof("login - unknown email: %q", email)
		// TODO: redirect to gateway with flash message
		return nil
	}

	if subtle.ConstantTimeCompare([]byte(user.Secret), []byte(pass)) == 0 {
		c.Log.Infof("login - invalid password for email %q", email)
		// TODO: redirect to gateway with flash message
		return nil
	}

	// User is now authed -> generate token.
	token, err := util.GenerateSecret(128)
	if err != nil {
		c.Log.Errorf("cannot generate secret: %q", err)
		// TODO: internal server error
		return nil
	}

	// Save token in cache.
	if err := cache.Set(token, struct{}{}, 168*time.Hour); err != nil {
		// TODO: can this happen at all?
		c.Log.Errorf("could not write token to cache: %q", err.Error())
	}

	// Save token in session.
	if err := c.Session.Set("access-token", token); err != nil {
		// TODO: should this be handled? I don't think this error can happen at all.
		c.Log.Errorf("could not write token to session: %q", err.Error())
		// TODO: internal server error
	}

	if remember {
		c.Session.SetDefaultExpiration()
	} else {
		c.Session.SetNoExpiration()
	}

	// Success -> redirect to home.
	return c.Redirect(App.Index)
}
