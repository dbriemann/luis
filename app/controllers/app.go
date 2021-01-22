package controllers

import (
	"bytes"
	"image"
	"luis/app/globals"
	"luis/app/models"
	"luis/app/store"
	"luis/app/util"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/revel/revel"
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
	user, ok := c.Args["user"].(models.User)
	if !ok {
		c.Log.Errorf("no user found in 'Args' for logged-in user")

		return c.RenderError(globals.ErrInternalServerError)
	}

	if err := user.FetchFiles(store.DB); err != nil {
		c.Log.Errorf("could not query files for user %q with error: %q", user.Email, err.Error())

		return c.RenderError(globals.ErrInternalServerError)
	}

	//put all files in viewargs
	c.ViewArgs["Files"] = user.Files

	return c.Render()
}

func (c App) File(id int64) revel.Result {
	return c.serveFile(id, false)
}

func (c App) Thumb(id int64) revel.Result {
	return c.serveFile(id, true)
}

func (c App) serveFile(id int64, thumb bool) revel.Result {
	user, ok := c.Args["user"].(models.User)
	if !ok {
		c.Log.Errorf("no user found in 'Args' for logged-in user")

		return c.RenderError(globals.ErrInternalServerError)
	}

	// Fetch file from DB.
	f, err := models.FileByID(store.DB, id)
	if err != nil {
		c.Log.Infof("could not retrieve File %d from db: %q", id, err.Error())

		return c.NotFound("file not found")
	}

	// If the files is not owned by the current user...
	if f.OwnerID != user.ID {
		// ... we check whether it's in one of the user's collections.
		// TODO

		return c.Forbidden("you have no access")
	}

	storagePath, found := revel.Config.String("storage.path")
	if !found {
		c.Log.Errorf("storage.path variable not set")

		return c.RenderError(globals.ErrInternalServerError)
	}

	var path string
	if thumb {
		path = filepath.Join(storagePath, f.Thumb)
	} else {
		path = filepath.Join(storagePath, f.Name)
	}

	// Serve the static file.
	return c.RenderFileName(path, revel.Inline)
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

	user, ok := c.Args["user"].(models.User)
	if !ok {
		c.Log.Errorf("no user found in 'Args' for logged-in user")

		return c.RenderError(globals.ErrInternalServerError)
	}

	// Extract file metadata.
	fname := c.Params.Files["file"][0].Filename

	storagePath, found := revel.Config.String("storage.path")
	if !found {
		c.Log.Errorf("storage.path variable not set")

		return c.RenderError(globals.ErrInternalServerError)
	}

	// Determine filetype and start according processing.
	// 1. Try to load as image.
	if img, err := imaging.Decode(bytes.NewReader(file)); err == nil {
		// It is an image -> make thumbnail and save image data.
		f, err := c.saveImage(storagePath, fname, img)
		if err != nil {
			c.Log.Errorf("error saving image data for file %q: %q", fname, err.Error())

			return c.RenderError(globals.ErrInternalServerError)
		}

		// Persist saved image in DB.
		// if err := gormdb.DB.Model(&user).Association("Files").Append(&f); err != nil {
		f.CreatedAt = time.Now().Unix()
		f.UpdatedAt = f.CreatedAt
		f.OwnerID = user.ID
		if err := f.Insert(store.DB); err != nil {
			c.Log.Errorf("could not persist file %q in DB with error %q", f.Name, err.Error())

			// In case of failure we have to clean up the orphaned image and thumb files.
			dstPathImage := filepath.Join(storagePath, f.Name)
			if err := os.Remove(dstPathImage); err != nil {
				c.Log.Errorf("could not clean up orphaned image at %q", dstPathImage)
				c.Log.Errorf("delete this image manually to avoid garbage")
			}

			dstPathThumb := filepath.Join(storagePath, f.Name)
			if err := os.Remove(dstPathThumb); err != nil {
				c.Log.Errorf("could not clean up orphaned thumb at %q", dstPathThumb)
				c.Log.Errorf("delete this image manually to avoid garbage")
			}

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

func (c App) saveImage(spath string, fname string, img image.Image) (models.File, error) {
	f := models.File{}

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
	finalThumbName := "thumb_" + finalName
	count := 1
	var dstPathImage string
	var dstPathThumb string
	for { // TODO: max iterations
		// Iterate until a free filename is found.
		dstPathImage = filepath.Join(spath, finalName)
		dstPathThumb = filepath.Join(spath, finalThumbName)
		_, ferr := os.Stat(dstPathImage)
		_, terr := os.Stat(dstPathThumb)

		if os.IsNotExist(ferr) && os.IsNotExist(terr) {
			// The file does not exist, get out of loop and use it.
			break
		}

		// If another error occurred we return it and fail.
		if ferr != nil && !os.IsNotExist(ferr) {
			c.Log.Debugf("image file stat error: %q", ferr.Error())
			return f, ferr
		}
		if terr != nil && !os.IsNotExist(terr) {
			c.Log.Debugf("thumb file stat error: %q", terr.Error())
			return f, terr
		}

		c.Log.Infof("filename %q or %q already exists", finalName, finalThumbName)

		// Find a new filename.
		count++
		prefix := strconv.Itoa(count) + "_"
		finalName = prefix + fname
		finalThumbName = "thumb_" + finalName
	}

	// Save image.
	err := imaging.Save(img, dstPathImage)
	if err != nil {
		return f, err
	}
	// Save thumb.
	err = imaging.Save(thumb, dstPathThumb)
	if err != nil {
		// If thumb cannot be saved we have to clean up the image.
		if err := os.Remove(dstPathImage); err != nil {
			c.Log.Errorf("could not clean up orphaned image at %q", dstPathImage)
			c.Log.Errorf("delete this image manually to avoid garbage")
		}

		return f, err
	}

	// CONTINUE
	f.Name = finalName
	f.Thumb = finalThumbName
	f.Type = models.FileTypeImage
	// TODO: where to handle collection?

	return f, nil
}

func (c App) ProfilePost() revel.Result {
	newemail := c.Params.Form.Get("email")
	newname := c.Params.Form.Get("name")

	user, ok := c.Args["user"].(models.User)
	if !ok {
		c.Log.Errorf("no user found in 'Args' for logged-in user")

		return c.RenderError(globals.ErrInternalServerError)
	}

	oldname := user.Name

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

		// if result := gormdb.DB.Save(&user); result.Error != nil {
		if err := user.Update(store.DB); err != nil {
			c.Log.Errorf("updating user (mail: %q, new mail: %q, name: %q, new name: %q) failed with error: %q", user.Email, newemail, oldname, newname, err.Error())

			return c.RenderError(globals.ErrInternalServerError)
		}

		c.Flash.Success("Thanks for your update.")
	}

	return c.Redirect(App.Profile)
}

func (c App) Profile() revel.Result {
	user, ok := c.Args["user"].(models.User)
	if !ok {
		c.Log.Errorf("no user found in 'Args' for logged-in user")

		return c.RenderError(globals.ErrInternalServerError)
	}

	c.ViewArgs["User"] = user

	return c.Render()
}
