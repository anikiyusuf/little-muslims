package utils 


import "os"


func IsDevMode() bool {
	return os.Getenv("APP_ENV") == "development"
}