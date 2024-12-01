package handlers

import (
	"EMtest/db"
	"EMtest/logger"
	"EMtest/models"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func addFilter(c *gin.Context, filterMap map[string]string, key string) {
	val, ok := c.GetQuery(key)
	if ok {
		filterMap[key] = val
	}
}

func processDBErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	logger.ERROR.Println(err.Error())
	status := http.StatusInternalServerError

	if errors.Is(err, gorm.ErrRecordNotFound) {
		status = http.StatusNotFound
	}

	c.JSON(status, gin.H{"error": err.Error()})
	return true
}

func getIDFromParam(c *gin.Context) (uint64, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		logger.ERROR.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID should be uint value"})
		return id, err
	}
	return id, nil
}

func getPageLimit(c *gin.Context) (int64, int64, error) {
	pageStr := c.DefaultQuery("page", "-1")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < -1 {
		log.Println("DeleteSong: ", err.Error())
		return 0, 0, err
	}

	limitStr := c.DefaultQuery("limit", "-1")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit < -1 {
		log.Println("DeleteSong: ", err.Error())
		return 0, 0, err
	}

	page--
	return page, limit, nil
}

func CreateSong(c *gin.Context) {
	var song models.Song

	if err := c.ShouldBindJSON(&song); err != nil {
		logger.ERROR.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if song.Name == "" || song.Group == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	logger.INFO.Println("CreateSong: creating song ", song)
	song, err := db.CreateSong(song)
	if processDBErr(c, err) {
		return
	}

	logger.INFO.Println("CreateSong: success")
	c.JSON(http.StatusCreated, song)
}

func CreateInfo(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		return
	}

	var songInfo models.SongInfo
	if err := c.ShouldBindJSON(&songInfo); err != nil {
		logger.ERROR.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "song not found in body"})
		return
	}
	if songInfo.Text == "" || songInfo.Link == "" || time.Time(songInfo.ReleaseDate).IsZero() {
		logger.INFO.Println("songinfo should have no null values")
		c.JSON(http.StatusBadRequest, gin.H{"error": "songinfo should have no null values"})
		return
	}
	songInfo.SongID = uint(id)

	songInfo, err = db.CreateInfo(songInfo)
	if processDBErr(c, err) {
		return
	}

	c.JSON(http.StatusCreated, songInfo)
}

func ChangeSongInfo(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		return
	}

	var songInfo models.SongInfo
	if err := c.ShouldBindJSON(&songInfo); err != nil {
		logger.ERROR.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "song not found in body"})
		return
	}

	updated, err := db.UpdateSongInfo(id, songInfo)
	if processDBErr(c, err) {
		return
	}

	c.JSON(http.StatusOK, updated)
}

func GetSong(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		return
	}
	logger.INFO.Println("GetSong: accessing db with parameters id#", id)

	song, err := db.GetSong(id)
	if processDBErr(c, err) {
		return
	}
	logger.INFO.Println("GetSong: success")

	c.JSON(http.StatusOK, song)
}

func GetSongs(c *gin.Context) {
	filterMap := make(map[string]string)

	addFilter(c, filterMap, "group")
	addFilter(c, filterMap, "song")
	addFilter(c, filterMap, "link")
	addFilter(c, filterMap, "release_date")
	addFilter(c, filterMap, "text")
	addFilter(c, filterMap, "like")

	page, limit, err := getPageLimit(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.INFO.Println("GetSongs: accessing db with parameters: filter", filterMap, "page=", page, "limit=", limit)
	songs, err := db.GetSongs(filterMap, page, limit)
	if processDBErr(c, err) {
		return
	}
	logger.INFO.Println("GetSongs: success")
	c.JSON(http.StatusOK, songs)
}

func GetSongText(c *gin.Context) {
	page, limit, err := getPageLimit(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := getIDFromParam(c)
	if err != nil {
		return
	}
	logger.INFO.Println("GetSongText: getting text of song#", id)

	text, err := db.GetSongText(id, page, limit)
	if processDBErr(c, err) {
		return
	}
	logger.INFO.Println("GetSongText: success")

	c.JSON(http.StatusOK, gin.H{"text": text})
}

func UpdateSong(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		return
	}

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		log.Println()
		c.JSON(http.StatusBadRequest, gin.H{"error": "song not found in body"})
		return
	}

	logger.INFO.Println("UpdateSong: updating song#", id)

	song, err = db.UpdateSong(id, song)
	if processDBErr(c, err) {
		return
	}

	logger.INFO.Println("UpdateSong: success")
	c.JSON(http.StatusOK, song)
}

func DeleteSong(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		return
	}

	logger.INFO.Println("DeleteSong: deleting song#", id)

	err = db.DeleteSong(id)
	if processDBErr(c, err) {
		return
	}
	logger.INFO.Println("DeleteSong: success")

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func DeleteSongInfo(c *gin.Context) {
	id, err := getIDFromParam(c)
	if err != nil {
		return
	}

	logger.INFO.Println("DeleteSongInfo: deleting song#", id)

	err = db.DeleteSongInfo(id)
	if processDBErr(c, err) {
		return
	}
	logger.INFO.Println("DeleteSongInfo: success")

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GetInfo(c *gin.Context) {
	groupName, ok := c.GetQuery("group")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No group name provided"})
		return
	}
	songName, ok := c.GetQuery("song")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No song name provided"})
		return
	}
	logger.INFO.Println("GetInfo: getting info of ", groupName, songName)

	info, err := db.GetSongInfo(groupName, songName)
	if processDBErr(c, err) {
		return
	}
	logger.INFO.Println("GetInfo: success")

	c.JSON(http.StatusOK, info)
}
