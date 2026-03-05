-- Current sql file was generated after introspecting the database
-- If you want to run this migration please uncomment this code before executing migrations
/*
CREATE TABLE `genres` (
	`id` int AUTO_INCREMENT NOT NULL,
	`code` varchar(50) NOT NULL,
	`name` varchar(100) NOT NULL,
	CONSTRAINT `genres_id` PRIMARY KEY(`id`),
	CONSTRAINT `code` UNIQUE(`code`)
);
--> statement-breakpoint
CREATE TABLE `movie_genres` (
	`movie_id` int NOT NULL,
	`genre_id` int NOT NULL,
	CONSTRAINT `movie_genres_movie_id_genre_id` PRIMARY KEY(`movie_id`,`genre_id`)
);
--> statement-breakpoint
CREATE TABLE `movies` (
	`id` int AUTO_INCREMENT NOT NULL,
	`tmdb_id` int NOT NULL,
	`title` varchar(200) NOT NULL,
	`overview` text NOT NULL,
	`release_date` date NOT NULL,
	`runtime_minutes` int NOT NULL,
	`poster_url` varchar(50),
	`series_id` int,
	`created_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
	`updated_at` timestamp DEFAULT (CURRENT_TIMESTAMP) ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT `movies_id` PRIMARY KEY(`id`),
	CONSTRAINT `tmdb_id` UNIQUE(`tmdb_id`),
	CONSTRAINT `movies_chk_1` CHECK((`runtime_minutes` > 0))
);
--> statement-breakpoint
CREATE TABLE `platforms` (
	`id` int AUTO_INCREMENT NOT NULL,
	`code` varchar(50) NOT NULL,
	`name` varchar(100) NOT NULL,
	CONSTRAINT `platforms_id` PRIMARY KEY(`id`),
	CONSTRAINT `code` UNIQUE(`code`)
);
--> statement-breakpoint
CREATE TABLE `refresh_tokens` (
	`id` int AUTO_INCREMENT NOT NULL,
	`user_id` int NOT NULL,
	`token_hash` varchar(255) NOT NULL,
	`expires_at` timestamp NOT NULL,
	`created_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
	`revoked_at` timestamp,
	CONSTRAINT `refresh_tokens_id` PRIMARY KEY(`id`),
	CONSTRAINT `token_hash` UNIQUE(`token_hash`)
);
--> statement-breakpoint
CREATE TABLE `reviews` (
	`id` int AUTO_INCREMENT NOT NULL,
	`user_id` int NOT NULL,
	`movie_id` int NOT NULL,
	`rating` decimal(2,1),
	`comment` text,
	`created_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
	`updated_at` timestamp DEFAULT (CURRENT_TIMESTAMP) ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT `reviews_id` PRIMARY KEY(`id`),
	CONSTRAINT `unique_user_movie_review` UNIQUE(`user_id`,`movie_id`),
	CONSTRAINT `reviews_chk_1` CHECK((`rating` between 0.1 and 5.0))
);
--> statement-breakpoint
CREATE TABLE `series` (
	`id` int AUTO_INCREMENT NOT NULL,
	`name` varchar(200) NOT NULL,
	`poster_url` varchar(50),
	`created_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
	`updated_at` timestamp DEFAULT (CURRENT_TIMESTAMP) ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT `series_id` PRIMARY KEY(`id`),
	CONSTRAINT `name` UNIQUE(`name`)
);
--> statement-breakpoint
CREATE TABLE `users` (
	`id` int AUTO_INCREMENT NOT NULL,
	`username` varchar(50) NOT NULL,
	`email` varchar(100) NOT NULL,
	`password_hash` varchar(255) NOT NULL,
	`is_active` tinyint(1) DEFAULT 1,
	`last_login_at` timestamp,
	`created_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
	`updated_at` timestamp DEFAULT (CURRENT_TIMESTAMP) ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT `users_id` PRIMARY KEY(`id`),
	CONSTRAINT `email` UNIQUE(`email`),
	CONSTRAINT `username` UNIQUE(`username`)
);
--> statement-breakpoint
CREATE TABLE `watch_history` (
	`id` int AUTO_INCREMENT NOT NULL,
	`user_id` int NOT NULL,
	`movie_id` int NOT NULL,
	`platform_id` int NOT NULL,
	`watched_date` date,
	`created_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
	CONSTRAINT `watch_history_id` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `watchlist` (
	`id` int AUTO_INCREMENT NOT NULL,
	`user_id` int NOT NULL,
	`movie_id` int NOT NULL,
	`priority` tinyint DEFAULT 1,
	`added_at` timestamp DEFAULT (CURRENT_TIMESTAMP),
	`updated_at` timestamp DEFAULT (CURRENT_TIMESTAMP) ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT `watchlist_id` PRIMARY KEY(`id`),
	CONSTRAINT `unique_user_movie` UNIQUE(`user_id`,`movie_id`),
	CONSTRAINT `watchlist_chk_1` CHECK((`priority` between 1 and 5))
);
--> statement-breakpoint
ALTER TABLE `movie_genres` ADD CONSTRAINT `movie_genres_ibfk_1` FOREIGN KEY (`movie_id`) REFERENCES `movies`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `movie_genres` ADD CONSTRAINT `movie_genres_ibfk_2` FOREIGN KEY (`genre_id`) REFERENCES `genres`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `movies` ADD CONSTRAINT `movies_ibfk_1` FOREIGN KEY (`series_id`) REFERENCES `series`(`id`) ON DELETE set null ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `refresh_tokens` ADD CONSTRAINT `refresh_tokens_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `reviews` ADD CONSTRAINT `reviews_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `reviews` ADD CONSTRAINT `reviews_ibfk_2` FOREIGN KEY (`movie_id`) REFERENCES `movies`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `watch_history` ADD CONSTRAINT `watch_history_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `watch_history` ADD CONSTRAINT `watch_history_ibfk_2` FOREIGN KEY (`movie_id`) REFERENCES `movies`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `watch_history` ADD CONSTRAINT `watch_history_ibfk_3` FOREIGN KEY (`platform_id`) REFERENCES `platforms`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `watchlist` ADD CONSTRAINT `watchlist_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE `watchlist` ADD CONSTRAINT `watchlist_ibfk_2` FOREIGN KEY (`movie_id`) REFERENCES `movies`(`id`) ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
CREATE INDEX `genre_id` ON `movie_genres` (`genre_id`);--> statement-breakpoint
CREATE INDEX `idx_release_date` ON `movies` (`release_date`);--> statement-breakpoint
CREATE INDEX `idx_series_id` ON `movies` (`series_id`);--> statement-breakpoint
CREATE INDEX `idx_expires_at` ON `refresh_tokens` (`expires_at`);--> statement-breakpoint
CREATE INDEX `idx_user_id` ON `refresh_tokens` (`user_id`);--> statement-breakpoint
CREATE INDEX `idx_created_at` ON `reviews` (`created_at`);--> statement-breakpoint
CREATE INDEX `movie_id` ON `reviews` (`movie_id`);--> statement-breakpoint
CREATE INDEX `idx_user_watched` ON `watch_history` (`user_id`,`watched_date`);--> statement-breakpoint
CREATE INDEX `idx_watched_date` ON `watch_history` (`watched_date`);--> statement-breakpoint
CREATE INDEX `movie_id` ON `watch_history` (`movie_id`);--> statement-breakpoint
CREATE INDEX `platform_id` ON `watch_history` (`platform_id`);--> statement-breakpoint
CREATE INDEX `idx_added_at` ON `watchlist` (`added_at`);--> statement-breakpoint
CREATE INDEX `movie_id` ON `watchlist` (`movie_id`);
*/