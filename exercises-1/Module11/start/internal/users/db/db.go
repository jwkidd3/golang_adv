package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var db *sql.DB

func initDB() {

	db = connectDB()

	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(4)
	db.SetConnMaxLifetime(time.Second * 15)
}

//GetDB returns the database handle
func Get() *sql.DB {
	if db == nil {
		initDB()
	}
	return db
}

func getDBConfig() (username string, password string,
	databasename string, databaseHost string) {
	dir, _ := os.Getwd()
	viper.SetConfigName("app")
	// Set the path to look for the configurations file
	viper.AddConfigPath(dir + "/../configs")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	databasename = viper.GetString("MYSQL_DATABASE")
	databaseHost = viper.GetString("MYSQL_SERVICE_HOST")
	username = viper.GetString("MYSQL_USERNAME")
	password = viper.GetString("MYSQL_PASSWORD")

	return
}

func createAndOpen(name string, dbURI string) (*sql.DB, error) {

	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		return nil, err
	}

	db, err = sql.Open("mysql", dbURI+name)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
						id int NOT NULL AUTO_INCREMENT,
						name varchar(100) NOT NULL,
						email varchar(100) NOT NULL,
						password varchar(100) NOT NULL,
						PRIMARY KEY (id)
					);`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

//ConnectDB function: Make database connection
func connectDB() *sql.DB {

	username, password, databasename, databaseHost := getDBConfig()

	//Define DB connection string
	dbURI := fmt.Sprintf("%s:%s@(%s)/", username, password, databaseHost)

	//connect to db URI
	db, err := createAndOpen(databasename, dbURI)
	if err != nil {
		fmt.Println("error", err)
		log.Fatalf(err.Error())
	}
	//ping the db cause it might be not really open
	err = db.Ping()

	if err != nil {
		fmt.Println("error", err)
		log.Fatalf(err.Error())
	}
	// close db when not in use
	// defer db.Close()
	log.Println("Successfully connected to db!")
	return db
}
