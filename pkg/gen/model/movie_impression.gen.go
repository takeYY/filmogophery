// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameMovieImpression = "movie_impression"

// MovieImpression mapped from table <movie_impression>
type MovieImpression struct {
	ID      int32   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	MovieID *int32  `gorm:"column:movie_id" json:"movie_id"`
	Status  bool    `gorm:"column:status;not null" json:"status"`
	Rating  *bool   `gorm:"column:rating" json:"rating"`
	Note    *string `gorm:"column:note" json:"note"`
	Movie   Movie   `gorm:"default:null;foreignKey:MovieID;references:ID" json:"movie"`
}

// TableName MovieImpression's table name
func (*MovieImpression) TableName() string {
	return TableNameMovieImpression
}
