// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameGenre = "genre"

// Genre mapped from table <genre>
type Genre struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	Code      string         `gorm:"column:code;not null" json:"code"`
	Name      string         `gorm:"column:name" json:"name"`
	Movies    []Movie        `gorm:"many2many:movie_genres" json:"movies"`
}

// TableName Genre's table name
func (*Genre) TableName() string {
	return TableNameGenre
}
