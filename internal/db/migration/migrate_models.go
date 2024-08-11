package migration

import (
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Name             string
		SelfIntroduction string
	}

	/*
		CREATE TABLE
		`poster` (
			`id` int AUTO_INCREMENT,
			`url` varchar(255) NOT NULL,
			PRIMARY KEY (`id`)
		);
	*/
	Poster struct {
		gorm.Model
		URL string `gorm:"column:url;not null"`
	}

	/*
		CREATE TABLE
		`movie_series` (
			`id` int AUTO_INCREMENT,
			`name` varchar(255) NOT NULL,
			`poster_id` int,
			PRIMARY KEY (`id`),
			FOREIGN KEY (`poster_id`) REFERENCES `poster` (`id`) ON DELETE SET NULL
		);
	*/
	MovieSeries struct {
		gorm.Model
		Name     string `gorm:"not null"`
		PosterID uint
		Poster   Poster
	}

	/*
		CREATE TABLE
		`genre` (
			`id` int AUTO_INCREMENT,
			`code` varchar(255) NOT NULL UNIQUE,
			`name` varchar(255),
			PRIMARY KEY (`id`)
		);
	*/
	Genre struct {
		gorm.Model
		Code   string  `gorm:"not null;unique"`
		Name   string  `gorm:"column:name"`
		Movies []Movie `gorm:"many2many:movie_genres"`
	}

	/*
		CREATE TABLE
		`movie` (
			`id` int AUTO_INCREMENT,
			`title` varchar(255) not null,
			`overview` text,
			`release_date` date not null,
			`run_time` int not null,
			`poster_id` int,
			`series_id` int,
			primary key (`id`),
			FOREIGN KEY (`poster_id`) REFERENCES `poster` (`id`) ON DELETE SET NULL,
			FOREIGN KEY (`series_id`) REFERENCES `movie_series` (`id`) ON DELETE SET NULL
		);
	*/
	Movie struct {
		gorm.Model
		Title       string      `gorm:"column:title;not null"`
		Overview    string      `gorm:"column:overview"`
		ReleaseDate time.Time   `gorm:"column:release_date;not null;"`
		RunTime     int32       `gorm:"column:run_time;not null;"`
		PosterID    uint        `gorm:"column:poster_id"`
		Poster      Poster      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		SeriesID    uint        `gorm:"column:series_id"`
		Series      MovieSeries `gorm:"foreignKey:SeriesID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Genres      []Genre     `gorm:"many2many:movie_genres"`
		TMDbID      int32       `gorm:"column:tmdb_id;not null"`
	}

	/*
		CREATE TABLE
		`watch_media` (
			`id` int AUTO_INCREMENT,
			`code` varchar(255) NOT NULL UNIQUE,
			`name` varchar(255),
			PRIMARY KEY (`id`)
		);
	*/
	WatchMedia struct {
		gorm.Model
		Code string `gorm:"not null;unique"`
		Name string
	}

	/*
		CREATE TABLE
		`movie_impression` (
			`id` int AUTO_INCREMENT,
			`movie_id` int,
			`status` int NOT NULL DEFAULT 0,
			`rating` int,
			`note` TEXT,
			PRIMARY KEY (`id`),
			FOREIGN KEY (`movie_id`) REFERENCES `movie` (`id`) ON DELETE CASCADE
		);
	*/
	MovieImpression struct {
		gorm.Model
		MovieID uint
		Movie   Movie
		Status  int32  `gorm:"column:status;not null"`
		Rating  int32  `gorm:"column:rating"`
		Note    string `gorm:"column:note"`
	}

	/*
		CREATE TABLE
		`movie_watch_record` (
			`id` int AUTO_INCREMENT,
			`watch_media_id` int,
			`watch_date` date,
			PRIMARY KEY (`id`),
			FOREIGN KEY (`watch_media_id`) REFERENCES `watch_media` (`id`) ON DELETE CASCADE
		);
	*/
	MovieWatchRecord struct {
		gorm.Model
		WatchMediaID uint
		WatchMedia   WatchMedia
		WatchDate    time.Time
	}
)
