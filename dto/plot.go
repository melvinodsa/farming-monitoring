package dto

import (
	"github.com/melvinodsa/farming-monitoring/db"
	"gorm.io/gorm"
)

//Plot dto for api
type Plot struct {
	ID uint
	//Name is the name of the plot
	Name string
	//Unit is the unit of plot area measurement
	Unit string
	//Location is the location at which the plot exists
	Location string
}

//ToPlot convert the from plot dto to db model
func (p Plot) ToPlot() db.Plot {
	return db.Plot{
		Model:    gorm.Model{ID: p.ID},
		Name:     p.Name,
		Unit:     p.Unit,
		Location: p.Location,
	}
}

//FromPlot convert the from db model to dto
func FromPlot(p db.Plot) Plot {
	return Plot{
		ID:       p.ID,
		Name:     p.Name,
		Unit:     p.Unit,
		Location: p.Location,
	}
}
