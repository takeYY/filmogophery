CREATE TABLE
    `poster` (
        `id` int AUTO_INCREMENT,
        `url` varchar(255) NOT NULL,
        PRIMARY KEY (`id`)
    );

CREATE TABLE
    `movie_series` (
        `id` int AUTO_INCREMENT,
        `name` varchar(255) NOT NULL,
        `poster_id` int,
        PRIMARY KEY (`id`),
        FOREIGN KEY (`poster_id`) REFERENCES `poster` (`id`) ON DELETE SET NULL
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
        `title` varchar(255) not null,
        `overview` text,
        `release_date` date not null,
        `run_time` int not null,
        `poster_id` int,
        `series_id` int,
        `tmdb_id` int,
        primary key (`id`),
        FOREIGN KEY (`poster_id`) REFERENCES `poster` (`id`) ON DELETE SET NULL,
        FOREIGN KEY (`series_id`) REFERENCES `movie_series` (`id`) ON DELETE SET NULL
    );

CREATE TABLE
    `movie_genres` (
        `id` int AUTO_INCREMENT,
        `movie_id` int not null,
        `genre_id` int not null,
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
        `movie_id` int,
        `status` tinyint (1) NOT NULL DEFAULT 0,
        `rating` tinyint (1),
        `note` TEXT,
        PRIMARY KEY (`id`),
        FOREIGN KEY (`movie_id`) REFERENCES `movie` (`id`) ON DELETE CASCADE
    );

CREATE TABLE
    `movie_watch_record` (
        `id` int AUTO_INCREMENT,
        `watch_media_id` int,
        `watch_date` date,
        PRIMARY KEY (`id`),
        FOREIGN KEY (`watch_media_id`) REFERENCES `watch_media` (`id`) ON DELETE CASCADE
    );