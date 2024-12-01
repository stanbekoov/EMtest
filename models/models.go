package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

const dateFormat = "02.01.2006"

type Date time.Time

type Song struct {
	ID    uint      `json:"id" gorm:"primary_key"`
	Name  string    `json:"song" gorm:"column:song"`
	Group string    `json:"group"`
	Info  *SongInfo `json:"info,omitempty" gorm:"foreignKey:SongID;OnDelete:CASCADE;"`
}

type SongInfo struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	SongID      uint   `json:"SongID" gorm:"unique;not null"`
	ReleaseDate Date   `json:"releaseDate" gorm:"type:date"`
	Link        string `json:"link"`
	Text        string `json:"text"`
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)

	t, err := time.Parse(dateFormat, s)
	if err != nil {
		return err
	}

	*d = Date(t)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format(dateFormat))
}

func (d Date) Format() string {
	return time.Time(d).Format(dateFormat)
}

func (d *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*d = Date(nullTime.Time)
	return err
}

func (d Date) Value() (driver.Value, error) {
	return time.Time(d), nil
}
