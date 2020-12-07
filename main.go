package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/farming-monitoring/config"
	"github.com/melvinodsa/farming-monitoring/db"
	"github.com/melvinodsa/farming-monitoring/routes"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := fiber.New()

	err := config.InitRootContext()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error while initializing the root context")
	}

	err = db.RunMigrations(config.NewContext())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("error while running db migrations")
	}

	routes.InitRoutes(app)

	log.Infof("started the server and listening at %s", config.GetConfig().Server.Port)
	app.Listen(":" + config.GetConfig().Server.Port)
}
