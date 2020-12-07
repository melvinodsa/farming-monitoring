package config

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//Context has the context with required utilities for a web request
type Context struct {
	//Log has the logger tagged with the user info
	Log *log.Logger
	//Db the db instance to communicate with the database
	Db *gorm.DB
}

var rootContext Context

//InitRootContext initializes the root context. It must be run after configs are loaded during the app boot
func InitRootContext() error {
	/*
	 * We will initialize the database instance
	 * We will initialize root context
	 */
	//initializing the database instance
	cfg := GetConfig()
	d, err := cfg.Db.ConnectWithRetry()
	if err != nil {
		log.Error("error while connecting to the database")
		return err
	}

	//initializing the root context
	rootContext.Db = d
	return nil
}

//NewContext returns a new context with database and other essentials utilities setup
func NewContext() Context {
	return Context{Log: log.New(), Db: rootContext.Db}
}

//AddContext middleware
func AddContext(c *fiber.Ctx) error {
	// Set some security headers
	c.Context().SetUserValue("farmingCtx", NewContext())

	// Go to next middleware:
	return c.Next()
}

//GetContext returns the config context
//when error happens, function writes the response. So no need of writing the response inside routes
func GetContext(c *fiber.Ctx) (*Context, error) {
	ctx := c.Context().UserValue("farmingCtx")
	if ctx == nil {
		log.Error("error while getting the app context. Context is nil")
		return nil, c.Status(500).JSON(Error("error while getting context"))
	}
	appCtx, ok := ctx.(Context)
	if !ok {
		log.Errorf("error while type converting the app context. Context type %q", reflect.TypeOf(ctx))
		return nil, c.Status(500).JSON(Error("error while getting context"))
	}
	return &appCtx, nil
}
