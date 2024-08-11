// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameMovieImpression = "movie_impression"

// MovieImpression mapped from table <movie_impression>
type MovieImpression struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	MovieID   int64          `gorm:"column:movie_id" json:"movie_id"`
	Status    int32          `gorm:"column:status;not null" json:"status"`
	Rating    int32          `gorm:"column:rating" json:"rating"`
	Note      string         `gorm:"column:note" json:"note"`
	Movie     Movie          `gorm:"foreignKey:ID" json:"movie"`
}

// TableName MovieImpression's table name
func (*MovieImpression) TableName() string {
	return TableNameMovieImpression
}
