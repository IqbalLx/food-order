package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env interface {
	Read(string) string
	ReadWithDefaultVal(envName string, defaultValue string) string
}

type env struct {
	envFile string
}

func NewEnv() *env {
	envFile := "./.env"
	_, err := os.Stat(envFile)

	if os.IsNotExist(err) {
		// Assume env variable injected directly from cli
		return &env{}
	} else if err != nil {
		// Handle other errors that may have occurred during the Stat operation
		log.Fatalln("Error checking file:", err)
	} else {
		godotenv.Load(envFile)
		return &env{envFile}
	}

	// default return, this line shouldn't ever be executed
	return &env{}
}


func (e *env) Read(envName string) string {
	return e.ReadWithDefaultVal(envName, "")
}

func (e *env) ReadWithDefaultVal(envName string, defaultValue string) string {
	value, isPresent := os.LookupEnv(envName)
	if !isPresent {
		if defaultValue != "" {
			return defaultValue
		}

		msg := fmt.Errorf("%s is not found in %s, please update", envName, e.envFile)
		panic(msg)
	}

	return value
}