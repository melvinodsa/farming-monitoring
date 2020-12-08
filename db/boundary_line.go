package db

import (
	"github.com/melvinodsa/farming-monitoring/config"
	"gorm.io/gorm"
)

//BoundaryLine indicates a boundary line of a plot
type BoundaryLine struct {
	gorm.Model
	//PlotID is the id of the plot to which the boundary line belongs to
	PlotID uint
	//X1 is x coordinate of the starting point of line
	X1 int
	//Y1 is y coordinate of the starting point of line
	Y1 int
	//X2 is x coordinate of the ending point of line
	X2 int
	//Y2 is y coordinate of the ending point of line
	Y2 int
}

//UpdateBoundaryLines deletes the old boundary lines and creates the new boundary lines
func UpdateBoundaryLines(ctx *config.Context, oldBoundaryLines, newBoundaryLines []BoundaryLine) ([]BoundaryLine, error) {
	/*
	 * First we will begin the transaction
	 * Then we will create the list of boundary line ids
	 * Then we will delete the boundary lines
	 * Then we will create/update the new boundary lines
	 */
	//creating the transaction
	tx := ctx.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, err
	}

	//build the list of ids of the boundary lines to be created
	ids := []uint{}
	for _, boundaryLine := range oldBoundaryLines {
		ids = append(ids, boundaryLine.ID)
	}

	//delete the boundary lines
	err := tx.Delete(&BoundaryLine{}, ids).Error
	if err != nil {
		tx.Rollback()
		ctx.Log.Errorf("error while deleting the boundary lines %+v", ids)
		return nil, err
	}

	//creating new boundary lines
	err = tx.Create(&newBoundaryLines).Error
	if err != nil {
		tx.Rollback()
		ctx.Log.Errorf("error while creating the new boundary lines")
		return nil, err
	}
	return newBoundaryLines, nil
}
