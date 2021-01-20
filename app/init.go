package app

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"luis/app/controllers"
	"luis/app/interceptors"
	"luis/app/models"
	"luis/app/store"
	"luis/app/util"
	"os"
	"path/filepath"
	"time"

	"github.com/revel/revel"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func ensurePaths() {
	revel.AppLog.Info("ensure paths exist")
	// Validate path for photo storage.
	photoPath, found := revel.Config.String("storage.path")
	if !found {
		revel.AppLog.Fatalf("no photo storage path set")
	}

	_, err := os.Stat(photoPath)
	if err != nil {
		if os.IsNotExist(err) {
			revel.AppLog.Fatalf("photo storage path %q does not exist", photoPath)
		} else {
			revel.AppLog.Fatalf("unexpected error: %q", err.Error())
		}
	}

	temp := []byte("temp")
	fname := filepath.Join(photoPath, "__access_test.___")
	if err := ioutil.WriteFile(fname, temp, 0600); err != nil {
		revel.AppLog.Fatalf("no write access in photo storage path %q", photoPath)
	}
	os.Remove(fname)
}

func ensureAdminAccess() {
	revel.AppLog.Info("ensure admin access")
	adminEmail, found := revel.Config.String("admin.email")
	if !found {
		revel.AppLog.Fatalf("no admin email set in ENV var %q", "LUIS_ADMIN_EMAIL")
	}

	revel.AppLog.Infof("LUIS_ADMIN_EMAIL set to %q", adminEmail)

	// Check if admin account already exists in DB. If not create it with random PW.
	admin, err := models.UserByEmail(store.DB, adminEmail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			secret, err := util.GenerateSecret(64)
			if err != nil {
				revel.AppLog.Fatalf("cannot generate high quality secret - No PRNG available")
			}

			admin.Email = adminEmail
			admin.Secret = secret
			admin.IsAdmin = true
			admin.CreatedAt = time.Now().Unix()
			admin.UpdatedAt = admin.CreatedAt

			// Create admin user (first start).
			if err := admin.Insert(store.DB); err != nil {
				revel.AppLog.Fatalf("cannot store admin user in DB: %q", err.Error())
			}
		} else {
			revel.AppLog.Fatalf("unexpected error: %q", err.Error())
		}
	}

	// needsUpdate := false
	//
	// If admin has no secret, generate one.
	if admin.Secret == "" {
		secret, err := util.GenerateSecret(64)
		if err != nil {
			revel.AppLog.Fatalf("cannot generate high quality secret - No PRNG available")
		}

		admin.Secret = secret
		admin.IsAdmin = true

		if err := admin.Update(store.DB); err != nil {
			revel.AppLog.Fatalf("cannot update user: %q with error: %q", admin.Email, err.Error())
		}
	}

	// TODO: if admin data is not complete show user view and force admin to enter name etc.

	// At this point we have a valid admin user stored in the database.
	revel.AppLog.Infof("admin secret is: %s", admin.Secret)
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		RememberRouteFilter,
		revel.ActionInvoker, // Invoke the action.
	}

	// Interceptors
	revel.InterceptFunc(interceptors.CheckAccess, revel.BEFORE, &controllers.App{})

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	revel.OnAppStart(store.InitDB, 0)
	revel.OnAppStart(ensurePaths, 1)
	revel.OnAppStart(ensureAdminAccess, 2)
}

var RememberRouteFilter = func(c *revel.Controller, fc []revel.Filter) {
	revel.AppLog.Debugf("route %s", c.Action)
	c.ViewArgs["Route"] = c.Action
	fc[0](c, fc[1:]) // Execute the next filter stage.
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
