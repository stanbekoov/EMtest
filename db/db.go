package db

import (
	"EMtest/logger"
	"EMtest/models"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	GORMlogger "gorm.io/gorm/logger"
)

var (
	db              *gorm.DB
	ErrSongNotFound = errors.New("song with these params doesn`t exist")
)

func Init() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: GORMlogger.New(
			log.New(logger.DEBUG.Writer(), "[GORM-DEBUG] ", log.Lshortfile|log.Ltime),
			GORMlogger.Config{
				SlowThreshold:             time.Second,
				Colorful:                  false,
				IgnoreRecordNotFoundError: true,
				ParameterizedQueries:      false,
				LogLevel:                  GORMlogger.Info,
			},
		),
	})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Song{}, &models.SongInfo{})
}

func addFiltersToQuery(old *gorm.DB, filterMap map[string]string) *gorm.DB {
	copy := old
	var word, symbol string
	if _, ok := filterMap["like"]; ok {
		word = "LIKE"
		symbol = "%"
	} else {
		word = "="
	}

	for k, v := range filterMap {
		if k == "like" {
			continue
		}

		condStr := fmt.Sprintf("\"%s\" %s ?", k, word)
		copy = copy.Where(condStr, symbol+v+symbol)
	}
	return copy
}

func CreateSong(song models.Song) (models.Song, error) {
	r := db.Model(&models.Song{}).Create(&song)
	return song, r.Error
}

func CreateInfo(info models.SongInfo) (models.SongInfo, error) {
	r := db.
		Model(&info).
		Clauses(clause.Returning{}).
		Create(&info)
	return info, r.Error
}

func GetSong(id uint64) (models.Song, error) {
	var song models.Song
	r := db.
		Model(&models.Song{}).
		Where(&models.Song{ID: uint(id)}).
		First(&song)

	return song, r.Error
}

func GetSongs(filterMap map[string]string, page, limit int64) ([]models.Song, error) {
	var songs []models.Song
	var offset int
	if page <= -1 && limit <= -1 {
		offset = -1
	} else {
		offset = int(page * limit)
	}

	query := db.
		Joins("Info").
		Model(&models.Song{}).
		Limit(int(limit)).
		Offset(offset)
	query = addFiltersToQuery(query, filterMap)
	r := query.Find(&songs)

	return songs, r.Error
}

func GetSongInfo(group, song string) (models.SongInfo, error) {
	var s models.Song
	r := db.
		Model(&models.Song{}).
		Where(models.Song{Group: group, Name: song}).
		Preload("Info").
		First(&s)
	if r.Error != nil {
		return models.SongInfo{}, r.Error
	}

	return *s.Info, nil
}

func GetSongText(id uint64, page, limit int64) (string, error) {
	if limit == -1 {
		limit = 4
	}

	var song models.Song
	r := db.
		Model(&models.Song{}).
		Preload("Info").
		First(&song, uint(id))
	if r.Error != nil {
		return "", r.Error
	}

	text := song.Info.Text
	if page == -1 {
		return text, nil
	}

	lines := strings.Split(text, "\n")
	left, right := page*limit, min(int64(len(lines)), (page+1)*limit)

	if left >= int64(len(lines)) || left < 0 {
		return "", nil
	}

	return strings.Join(lines[left:right], "\n"), nil
}

func DeleteSong(id uint64) error {
	return db.Delete(&models.Song{ID: uint(id)}).Error
}

func DeleteSongInfo(id uint64) error {
	return db.Model(&models.SongInfo{}).Where(&models.SongInfo{SongID: uint(id)}).Delete(&models.SongInfo{}).Error
}

func UpdateSong(id uint64, newSong models.Song) (models.Song, error) {
	var updated models.Song
	r := db.
		Model(&updated).
		Clauses(clause.Returning{}).
		Where(&models.Song{ID: uint(id)}).
		Updates(&newSong)
	return updated, r.Error
}

func UpdateSongInfo(id uint64, newInfo models.SongInfo) (models.SongInfo, error) {
	var updated models.SongInfo
	r := db.
		Model(&updated).
		Clauses(clause.Returning{}).
		Where(&models.SongInfo{SongID: uint(id)}).
		Updates(&newInfo)
	return updated, r.Error
}
