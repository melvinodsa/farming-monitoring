package config

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Database is the database connection config
type Database struct {
	//Name of the database
	Name string
	//Username to access the database
	Username string
	//Password of the username to access the database
	Password string
	//Host at which database is avaiable
	Host string
	//Port at which database is avaiable on the given host
	Port string
}

//Connect connects to database. Return db instance if sucessfully connected. Else return error
func (d Database) Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", d.Host, d.Username, d.Password, d.Name, d.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

//ConnectWithRetry connects to database with 5 retries.
func (d Database) ConnectWithRetry() (*gorm.DB, error) {
	retry := 5
	for retry > 0 {
		d, err := d.Connect()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Errorf("error while connecting to database. Retrying in 3 second. retry count %d", retry)
			retry--
			time.Sleep(3 * time.Second)
			continue
		}
		return d, nil
	}
	return nil, fmt.Errorf("couldn't connect to database after %d retries", 5)
}
