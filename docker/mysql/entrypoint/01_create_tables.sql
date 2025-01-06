CREATE TABLE
    `movie_series` (
        `id` int AUTO_INCREMENT,
        `name` varchar(255) NOT NULL,
        `poster_url` varchar(255),
        PRIMARY KEY (`id`)
    );

CREATE TABLE
    `genre` (
        `id` int AUTO_INCREMENT,
        `code` varchar(255) NOT NULL UNIQUE,
        `name` varchar(255),
        PRIMARY KEY (`id`)
    );

CREATE TABLE
    `movie` (
        `id` int AUTO_INCREMENT,
        `title` varchar(255) NOT NULL,
        `overview` text NOT NULL,
        `release_date` date NOT NULL,
        `run_time` int NOT NULL,
        `poster_url` varchar(255),
        `series_id` int,
        `tmdb_id` int NOT NULL,
        primary key (`id`),
        FOREIGN KEY (`series_id`) REFERENCES `movie_series` (`id`) ON DELETE SET NULL
    );

CREATE TABLE
    `movie_genres` (
        `id` int AUTO_INCREMENT,
        `movie_id` int NOT NULL,
        `genre_id` int NOT NULL,
        primary key (`id`),
        FOREIGN KEY (`movie_id`) REFERENCES `movie` (`id`) on delete cascade,
        FOREIGN KEY (`genre_id`) REFERENCES `genre` (`id`) on delete cascade
    );

CREATE TABLE
    `watch_media` (
        `id` int AUTO_INCREMENT,
        `code` varchar(255) NOT NULL UNIQUE,
        `name` varchar(255),
        PRIMARY KEY (`id`)
    );

CREATE TABLE
    `movie_impression` (
        `id` int AUTO_INCREMENT,
        `movie_id` int NOT NULL,
        `status` tinyint (1) NOT NULL DEFAULT 0,
        `rating` float (2, 1),
        `note` TEXT,
        PRIMARY KEY (`id`),
        FOREIGN KEY (`movie_id`) REFERENCES `movie` (`id`) ON DELETE CASCADE
    );

CREATE TABLE
    `movie_watch_record` (
        `id` int AUTO_INCREMENT,
        `movie_impression_id` int NOT NULL,
        `watch_media_id` int NOT NULL,
        `watch_date` date NOT NULL,
        PRIMARY KEY (`id`),
        FOREIGN KEY (`movie_impression_id`) REFERENCES `movie_impression` (`id`) ON DELETE CASCADE,
        FOREIGN KEY (`watch_media_id`) REFERENCES `watch_media` (`id`) ON DELETE CASCADE
    );