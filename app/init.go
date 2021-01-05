package app

import (
	"errors"
	"luis/app/controllers"
	"luis/app/gormdb"
	"luis/app/interceptors"
	"luis/app/models"
	"luis/app/util"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func ensureAdminAccess() {
	revel.AppLog.Warn("ensureAdminAccess()")
	adminEmail, found := revel.Config.String("admin.email")
	if !found {
		revel.AppLog.Fatalf("no admin email set in ENV var %q", "LUIS_ADMIN_EMAIL")
	}

	revel.AppLog.Infof("LUIS_ADMIN_EMAIL set to %q", adminEmail)

	// Check if admin account already exists in DB. If not create it with random PW.
	admin := models.User{
		Email: adminEmail,
	}
	result := gormdb.DB.Take(&admin)

	needsUpdate := false

	// If admin has no secret, generate one.
	if admin.Secret == "" {
		// Create initial admin user secret.
		secret, err := util.GenerateSecret(64)
		if err != nil {
			revel.AppLog.Fatalf("cannot generate high quality secret - No PRNG available")
		}

		admin.Secret = secret
		needsUpdate = true
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Admin user didn't exist in DB. Create.
		result := gormdb.DB.Create(&admin)
		if result.Error != nil {
			revel.AppLog.Fatalf("could not create admin user: %q", result.Error.Error())
		}
	} else if needsUpdate {
		// Update admin with new password.
		result := gormdb.DB.Save(&admin)
		if result.Error != nil {
			revel.AppLog.Fatalf("could not update admin user: %q", result.Error.Error())
		}
	}

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
		revel.ActionInvoker,           // Invoke the action.
	}

	// Interceptors
	revel.InterceptFunc(interceptors.CheckAccess, revel.BEFORE, &controllers.App{})

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	revel.OnAppStart(gormdb.InitDB, 0)
	revel.OnAppStart(ensureAdminAccess, 1)
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
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
