package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/melvinodsa/farming-monitoring/config"
)

//InitRoutes will initialize the routes for the app
func InitRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1, https://farming.melvindavis/me",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(config.AddContext)
	initPlotRoutes(app)
}

func initPlotRoutes(app *fiber.App) {
	plot := app.Group("/plot", PlotMiddleware)
	plot.Get("", GetPlots)
	plot.Put("", CreatePlot)
	plot.Post("", UpdatePlot)
	plot.Get("/:id", GetPlot)
	plot.Delete("/:id", DeletePlot)
}
