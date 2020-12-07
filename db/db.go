package db

import "github.com/melvinodsa/farming-monitoring/config"

//RunMigrations run the migrations required for the app database functioning
func RunMigrations(c config.Context) error {
	return c.Db.AutoMigrate(&Plot{})
}
