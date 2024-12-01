package main

import (
	"EMtest/db"
	"EMtest/handlers"
	"EMtest/logger"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	infoFile  *os.File
	debugFile *os.File
	errorFile *os.File
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	infoFile = logger.SetOutputByFilename(logger.INFO, os.Getenv("INFO_LOG_FILE"))
	debugFile = logger.SetOutputByFilename(logger.DEBUG, os.Getenv("DEBUG_LOG_FILE"))
	errorFile = logger.SetOutputByFilename(logger.DEBUG, os.Getenv("ERROR_LOG_FILE"))

	logger.INFO.Println("info logs output set to ", logger.INFO.Writer())

	db.Init()
}

func main() {
	if infoFile != nil {
		defer infoFile.Close()
	}
	if debugFile != nil {
		defer debugFile.Close()
	}
	if errorFile != nil {
		defer errorFile.Close()
	}

	gin.DefaultWriter = logger.INFO.Writer()
	gin.DefaultErrorWriter = logger.ERROR.Writer()
	router := gin.Default()

	router.GET("/info", handlers.GetInfo)
	router.GET("/songs", handlers.GetSongs)
	router.GET("/songs/:id", handlers.GetSong)
	router.GET("/songs/:id/text", handlers.GetSongText)
	router.POST("/songs", handlers.CreateSong)
	router.POST("/songs/:id/info", handlers.CreateInfo)
	router.DELETE("/songs/:id", handlers.DeleteSong)
	router.DELETE("/songs/:id/info", handlers.DeleteSongInfo)
	router.PATCH("/songs/:id", handlers.UpdateSong)
	router.PATCH("/songs/:id/info", handlers.ChangeSongInfo)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("you should specify port in .env file")
	}
	router.Run(":" + port)
}
