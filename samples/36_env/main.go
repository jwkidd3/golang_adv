package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golobby/dotenv"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	pl = fmt.Println
	lf = log.Fatalln

	envPath = ".env"
)

// Use default `os` package
func useOSPackage() {
	// Set ENV values
	svNameTitle := "SERVICE_NAME"
	os.Setenv(svNameTitle, "example")

	// Get values from an ENV
	svName := os.Getenv(svNameTitle)
	pl("Service name :", svName)

	// Unset ENV values
	os.Unsetenv(svNameTitle)
	pl("After unset, Service name :", os.Getenv(svNameTitle))

	// Check if an ENV value exists or not
	dbUrl, ok := os.LookupEnv("DB_URL")
	if !ok {
		pl("DB URL string does not exist")
	} else {
		pl("DB URL :", dbUrl)
	}
}

func useGoDotEnv() {
	/*
		Support file formats: .env, yaml
		Install lib
			go get github.com/joho/godotenv
	*/

	// godotenv.Load() // default load in the current directory
	err := godotenv.Load(envPath)
	if err != nil {
		lf("Error loading .env file")
	}

	pl("WAKANDA :", os.Getenv("WAKANDA"))
}

func useDotEnv() {
	/*
		Support file format: .env, yaml
		A lightweight package for loading OS environment variables into structs for Go projects

		Install lib
			go get github.com/golobby/dotenv
	*/
	type Config struct {
		Debug bool `env:"DEBUG"`
		App   struct {
			Name string `env:"APP_NAME"`
			Port string `env:"APP_PORT"`
		}
	}

	// Read `.env` file
	file, err := os.Open(envPath)
	if err != nil {
		lf("Error loading .env file")
	}

	var config Config
	if err = dotenv.NewDecoder(file).Decode(&config); err != nil {
		lf("Can not map ENV variables to config struct")
	}

	pl("APP_NAME :", config.App.Name)
}

func viperLoadFromEnv() {
	// Set env file and path
	viper.SetConfigFile(envPath)

	// Find and read that env file
	err := viper.ReadInConfig()
	if err != nil {
		lf("Error reading the env file :", err)
	}
}

func viperLoadFromYaml() {
	// Set name of config file (without file extension)
	viper.SetConfigName("config")

	// Look for file in current directory
	viper.AddConfigPath(".")

	// Find and read that config file
	if err := viper.ReadInConfig(); err != nil {
		lf("Error reading the config file :", err)
	}
}

func useViper() {
	/*	⭐️⭐️⭐️⭐️⭐️
		(common, recommend from the golang community)
		Support file format: .env, yaml, json, toml, hcl, Java properties

		Install lib
			go get github.com/spf13/viper
	*/
	viperLoadFromYaml()

	// viper.Get() returns an empty interface{}
	// to get the underlying type of the key,
	// we have to do the type assertion, we know the underlying value is string
	// if we type assert to other type it will throw an error
	value1, ok := viper.Get("I_AM").(string)

	// If the type is a string then ok will be true
	// ok will make sure the program not break
	if !ok {
		lf("Invalid type assertion")
	}

	// Get config values
	pl("config.yaml - I_AM :", value1)

	viperLoadFromEnv()
	value2, ok := viper.Get("DB_NAME").(string)
	if !ok {
		lf("Invalid type assertion")
	}
	pl(".env - DB_NAME :", value2)
}

func main() {
	pl("\n----- os -----")
	useOSPackage()

	pl("\n----- godotenv -----")
	useGoDotEnv()

	pl("\n----- dotenv -----")
	useDotEnv()

	pl("\n----- viper -----")
	useViper()
}
