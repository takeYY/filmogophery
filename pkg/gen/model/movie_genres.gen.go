// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameMovieGenres = "movie_genres"

// MovieGenres mapped from table <movie_genres>
type MovieGenres struct {
	MovieID int64 `gorm:"column:movie_id;primaryKey" json:"movie_id"`
	GenreID int64 `gorm:"column:genre_id;primaryKey" json:"genre_id"`
}

// TableName MovieGenres's table name
func (*MovieGenres) TableName() string {
	return TableNameMovieGenres
}
