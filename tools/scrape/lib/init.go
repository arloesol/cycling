package lib

import (
	"log"
	"os"
)

func init() {
	// logging system
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	LogInfo = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogWarning = log.New(file, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	LogError = log.New(file, " ERR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// data directory
	basedir := os.Getenv("GITDIR")
	if basedir == "" {
		LogError.Panicln("GITDIR env variable not set")
	}
	datadir = basedir + "/data"
}
