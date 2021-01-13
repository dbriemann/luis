package controllers

import (
	"bytes"
	"errors"
	"image"
	"luis/app/globals"
	"luis/app/gormdb"
	"luis/app/models"
	"luis/app/util"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/revel/revel"
	"gorm.io/gorm"
)

const (
	_      = iota
	KB int = 1 << (10 * iota)
	MB
	GB
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Upload() revel.Result {
	return c.Render()
}

func (c App) UploadPost(file []byte) revel.Result {
	c.Validation.Required(file)
	c.Validation.MaxSize(file, 1*GB)

	if c.Validation.HasErrors() {
		// TODO handle errors in frontend
		c.Validation.Keep()
		c.FlashParams()

		return c.Redirect(App.Upload)
	}

	// Extract file metadata.
	fname := c.Params.Files["file"][0].Filename

	// Determine filetype and start according processing.
	// 1. Try to load as image.
	img, err := imaging.Decode(bytes.NewReader(file))
	if err == nil {
		// It is an image -> make thumbnail and save image data.
		if err := c.saveImage(fname, img); err != nil {
			c.Log.Errorf("error saving image data for file %q: %q", fname, err.Error())

			return c.RenderError(globals.ErrInternalServerError)
		}

		return c.RenderText("OK")
	} else {
		c.Log.Debugf("file %q is not an image", fname)
		// TODO remove error and continue to handle other filetypes
		c.Response.Status = http.StatusBadRequest

		return c.RenderText("file format not supported")
	}

	// TODO: if video use ffmpeg to make thumbnail

	return c.RenderText("OK")
}

func (c App) saveImage(fname string, img image.Image) error {
	storagePath, found := revel.Config.String("storage.path")
	if !found {
		return globals.ErrStorageUnavailable
	}

	// Make thumbnail.
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	if w > h {
		w = globals.ThumbnailSizePixels
		h = 0
	} else {
		w = 0
		h = globals.ThumbnailSizePixels
	}
	thumb := imaging.Resize(img, w, h, imaging.Lanczos)

	finalName := fname
	count := 1
	var dstPathImage string
	var dstPathThumb string
	for {
		// Iterate until a free filename is found.
		dstPathImage = filepath.Join(storagePath, finalName)
		dstPathThumb = filepath.Join(storagePath, "thumb_"+finalName)
		_, ferr := os.Stat(dstPathImage)
		_, terr := os.Stat(dstPathThumb)

		if os.IsNotExist(ferr) && os.IsNotExist(terr) {
			// The file does not exist, get out of loop and use it.
			break
		}

		// If another error occured we return it and fail.
		if ferr != nil {
			return ferr
		}
		if terr != nil {
			return terr
		}

		c.Log.Infof("filename %q already exists")

		// Find a new filename.
		count++
		prefix := strconv.Itoa(count) + "_"
		finalName = prefix + fname
	}

	// Save image.
	err := imaging.Save(img, dstPathImage)
	if err != nil {
		return err
	}
	// Save thumb.
	err = imaging.Save(thumb, dstPathThumb)
	if err != nil {
		// TODO delete image.
		return err
	}

	// TODO handle metadata.
	// TODO persist in DB.

	return nil
}

func (c App) ProfilePost() revel.Result {
	newemail := c.Params.Form.Get("email")
	newname := c.Params.Form.Get("name")

	// TODO: assert after ok check
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
