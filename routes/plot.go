package routes

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/melvinodsa/farming-monitoring/config"
	"github.com/melvinodsa/farming-monitoring/db"
	"github.com/melvinodsa/farming-monitoring/dto"
	log "github.com/sirupsen/logrus"
)

const plotKey = "plot"

//PlotMiddleware parses the plot dto to db object and put it in the context
func PlotMiddleware(c *fiber.Ctx) error {
	/*
	 * We will get the context
	 * We will check if the request is put/post
	 * We will parse the plot dto
	 * extra check for post - update request
	 * Convert to db model
	 * Save it in the context
	 */
	//getting the context
	ctx, err := config.GetContext(c)
	if err != nil {
		//error while getting the context
		return err
	}

	//checking if the request is put/post
	requestType := string(c.Request().Header.Method())
	if requestType != fiber.MethodPut && requestType != fiber.MethodPost {
		return c.Next()
	}

	//parsing the dto object
	plot := dto.Plot{}
	err = c.BodyParser(&plot)
	if err != nil {
		//error while parsing the body
		ctx.Log.WithFields(log.Fields{"error": err}).Error("error while parsing the request body")
		return c.Status(400).JSON(config.Error("bad request. check the payload"))
	}

	if requestType == fiber.MethodPost && plot.ID == 0 {
		//request id missing
		ctx.Log.WithFields(log.Fields{"error": errors.New("plot id missing in the payload")}).Error("plot id missing in json payload")
		return c.Status(400).JSON(config.Error("bad request. check the payload. id missing"))
	}

	//we will convert the dto to db model
	dbPlot := plot.ToPlot()
	c.Context().SetUserValue(plotKey, dbPlot)

	return c.Next()
}

//GetPlots return the list of plot the user has access to
func GetPlots(c *fiber.Ctx) error {
	/*
	 * We will get the context
	 * We will get all the plots from db
	 * convert the db instance to dto
	 * Return the response
	 */
	//getting the context
	ctx, err := config.GetContext(c)
	if err != nil {
		//error while getting the context
		return err
	}

	//finding all the plots from the database
	plots, err := db.FindAllPlots(ctx)
	if err != nil {
		ctx.Log.WithFields(log.Fields{"error": err}).Error("couldn't find the plots from db")
		return c.Status(500).JSON(config.Error("couldn' find the plots"))
	}

	//converting the db instance to dto
	resp := []dto.Plot{}
	for _, p := range plots {
		resp = append(resp, dto.FromPlot(p))
	}

	//giving out the response
	return c.JSON(config.Success(resp, "plots"))
}

//CreatePlot will create the plot
func CreatePlot(c *fiber.Ctx) error {
	/*
	 * Getting the context
	 * Fetching the request payload
	 * Creating the plot in the database
	 * Returning the response
	 */
	//getting the context
	ctx, err := config.GetContext(c)
	if err != nil {
		//error while getting the context
		return err
	}

	//parsing the request body
	plotCtx := c.Context().UserValue(plotKey)
	plot, ok := plotCtx.(db.Plot)
	if !ok {
		//error while type converting the db plot instance
		ctx.Log.WithFields(log.Fields{"error": err}).Error("error while processing the plot info")
		return c.Status(500).JSON(config.Error("error while procssing the plot info"))
	}

	//creating the db entry
	err = (&plot).Create(ctx)
	if err != nil {
		//error while creating the database entry
		ctx.Log.WithFields(log.Fields{"error": err}).Error("error while creating the database entry for the plot")
		return c.Status(500).JSON(config.Error("couldn't save the plot into the db"))
	}

	//returning the response
	return c.JSON(config.Success(dto.FromPlot(plot), "plot"))
}

//GetPlot will return the plot corresponding to the given id
func GetPlot(c *fiber.Ctx) error {
	/*
	 * Getting the context
	 * getting the plot id
	 * finding the plot info from db
	 * returning the response
	 */
	//getting the context
	ctx, err := config.GetContext(c)
	if err != nil {
		//error while getting the context
		return err
	}

	//getting the params
	id := c.Params("id")
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		//error while parsing the plot id
		ctx.Log.WithFields(log.Fields{"error": err}).Errorf("given plot id %s not a valid uint", id)
		return c.Status(422).JSON(config.Error("plot id not a valid integer"))
	}

	//finding the plot from db
	plot := &db.Plot{}
	err = plot.FindByID(ctx, uint(parsedId))
	if err != nil {
		//couldn't find the plot for the given id
		ctx.Log.WithFields(log.Fields{"error": err}).Errorf("given plot id %d doesn't exist", parsedId)
		return c.Status(404).JSON(config.Error("plot not found"))
	}

	//return the response
	return c.JSON(config.Success(dto.FromPlot(*plot), "plot"))
}

//UpdatePlot will update the plot corresponding to the given id
func UpdatePlot(c *fiber.Ctx) error {
	/*
	 * Getting the context
	 * Fetching the request payload
	 * Updating the plot in the database
	 * Returning the response
	 */
	//getting the context
	ctx, err := config.GetContext(c)
	if err != nil {
		//error while getting the context
		return err
	}

	//parsing the request body
	plotCtx := c.Context().UserValue(plotKey)
	plot, ok := plotCtx.(db.Plot)
	if !ok {
		//error while type converting the db plot instance
		ctx.Log.WithFields(log.Fields{"error": err}).Error("error while processing the plot info")
		return c.Status(500).JSON(config.Error("error while procssing the plot info"))
	}

	//updating the db entry
	err = (&plot).Update(ctx)
	if err != nil {
		//error while updating the database entry
		ctx.Log.WithFields(log.Fields{"error": err}).Error("error while updating the database entry for the plot")
		return c.Status(500).JSON(config.Error("couldn't update the plot in the db"))
	}

	//returning the response
	return c.JSON(config.Success(dto.FromPlot(plot), "plot"))
}

//DeletePlot will delete the plot corresponding to the given id
func DeletePlot(c *fiber.Ctx) error {
	/*
	 * Getting the context
	 * getting the plot id
	 * deleting the plot info from db
	 * returning the response
	 */
	//getting the context
	ctx, err := config.GetContext(c)
	if err != nil {
		//error while getting the context
		return err
	}

	//getting the params
	id := c.Params("id")
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		//error while parsing the plot id
		ctx.Log.WithFields(log.Fields{"error": err}).Errorf("given plot id %s not a valid uint", id)
		return c.Status(422).JSON(config.Error("plot id not a valid integer"))
	}

	//finding the plot from db
	plot := &db.Plot{}
	err = plot.Delete(ctx, uint(parsedId))
	if err != nil {
		//couldn't find the plot for the given id
		ctx.Log.WithFields(log.Fields{"error": err}).Errorf("given plot id %d doesn't exist", parsedId)
		return c.Status(404).JSON(config.Error("plot not found"))
	}

	//return the response
	return c.JSON(config.Success(nil, "plot deleted"))
}
