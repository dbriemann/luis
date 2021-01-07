package controllers

import (
	"errors"
	"luis/app/globals"
	"luis/app/gormdb"
	"luis/app/models"
	"luis/app/util"

	"github.com/revel/revel"
	"gorm.io/gorm"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) ProfilePost() revel.Result {
	newemail := c.Params.Form.Get("email")
	newname := c.Params.Form.Get("name")

	email, ok := c.Args["email"].(string)
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

	oldname := user.Name

	// c.ViewArgs["User"] = user
	// TODO: switch to activeEmail and Email for activation of new email?

	// Compare old and new data.
	if newemail != user.Email || newname != user.Name {
		// Something changed -> update.
		if !util.IsEMailValid(newemail) {
			// Invalid data.
			c.Flash.Error(globals.ErrInvalidEmail.Error())

			return c.Redirect(App.Profile)
		}

		user.Email = newemail
		user.Name = newname

		if result := gormdb.DB.Save(&user); result.Error != nil {
			c.Log.Errorf("updating user (mail: %q, new mail: %q, name: %q, new name: %q) failed: %q", email, newemail, oldname, newname)

			return c.RenderError(globals.ErrInternalServerError)
		}

		c.Flash.Success("Thanks for your update.")
	}

	return c.Redirect(App.Profile)
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
