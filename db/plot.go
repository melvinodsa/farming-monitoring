package db

import (
	"github.com/melvinodsa/farming-monitoring/config"
	"gorm.io/gorm"
)

//Plot has the db model required for storing plot information
type Plot struct {
	gorm.Model
	//Name is the name of the plot
	Name string
	//Unit is the unit of plot area measurement
	Unit string
	//Location is the location at which the plot exists
	Location string
}

//FindAllPlots gets all the plots in the db
func FindAllPlots(ctx *config.Context) ([]Plot, error) {
	plots := []Plot{}
	err := ctx.Db.Find(&plots).Error
	return plots, err
}

//Create will save the plot into database
func (p *Plot) Create(ctx *config.Context) error {
	return ctx.Db.Create(p).Error
}

//FindByID finds the plot for the given id
func (p *Plot) FindByID(ctx *config.Context, id uint) error {
	return ctx.Db.First(p, "id = ?", id).Error
}

//Update will update the name/unit/location of a plot
func (p *Plot) Update(ctx *config.Context) error {
	pU := Plot{Model: gorm.Model{ID: p.ID}}
	err := ctx.Db.Model(&pU).Updates(map[string]interface{}{"name": p.Name, "unit": p.Unit, "location": p.Location}).Error
	*p = pU
	return err
}

//Delete will soft delete the plot from db the given id
func (p *Plot) Delete(ctx *config.Context, id uint) error {
	err := p.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return ctx.Db.Delete(&Plot{}, id).Error
}
