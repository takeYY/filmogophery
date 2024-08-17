// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameMovieWatchRecord = "movie_watch_record"

// MovieWatchRecord mapped from table <movie_watch_record>
type MovieWatchRecord struct {
	ID                int32      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	MovieImpressionID int32      `gorm:"column:movie_impression_id;not null" json:"movie_impression_id"`
	WatchMediaID      int32      `gorm:"column:watch_media_id;not null" json:"watch_media_id"`
	WatchDate         time.Time  `gorm:"column:watch_date;not null" json:"watch_date"`
	WatchMedia        WatchMedia `gorm:"foreignKey:WatchMediaID;references:ID" json:"watch_media"`
}

// TableName MovieWatchRecord's table name
func (*MovieWatchRecord) TableName() string {
	return TableNameMovieWatchRecord
}