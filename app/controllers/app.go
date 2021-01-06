package controllers

import (
	"errors"
	"luis/app/globals"
	"luis/app/gormdb"
	"luis/app/models"

	"github.com/revel/revel"
	"gorm.io/gorm"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Profile() revel.Result {
	arg := c.Args["email"]
	email, ok := arg.(string)

	if !ok {
		c.Log.Errorf("no email found in 'Args' for logged-in user")

		return c.RenderError(globals.ErrInternalServerError)
	}

	// Fetch user data from DB.
	user := models.User{
		Email: email,
	}

	result := gormdb.DB.Take(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Log.Errorf("no user found in DB for logged-in user")

		return c.RenderError(globals.ErrInternalServerError)
	}

	c.ViewArgs["User"] = user

	return c.Render()
}
