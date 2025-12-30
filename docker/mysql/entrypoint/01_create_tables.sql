USE db4dev;

-- ユーザーテーブル
CREATE TABLE
    `users` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `username` VARCHAR(50) NOT NULL UNIQUE,
        `email` VARCHAR(100) NOT NULL UNIQUE,
        `password_hash` VARCHAR(255) NOT NULL,
        `is_active` BOOLEAN DEFAULT TRUE,
        `last_login_at` TIMESTAMP NULL,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- リフレッシュトークンテーブル
CREATE TABLE
    `refresh_tokens` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `user_id` INT NOT NULL,
        `token_hash` VARCHAR(255) NOT NULL UNIQUE,
        `expires_at` TIMESTAMP NOT NULL,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `revoked_at` TIMESTAMP NULL,
        FOREIGN KEY (user_id) REFERENCES `users` (id) ON DELETE CASCADE,
        INDEX idx_user_id (user_id),
        INDEX idx_expires_at (expires_at)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- ジャンルテーブル
CREATE TABLE
    `genres` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `code` VARCHAR(50) NOT NULL UNIQUE,
        `name` VARCHAR(100) NOT NULL
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- シリーズテーブル
CREATE TABLE
    `series` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `name` VARCHAR(200) NOT NULL UNIQUE,
        `poster_url` VARCHAR(50),
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 映画テーブル
CREATE TABLE
    `movies` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `tmdb_id` INT NOT NULL UNIQUE,
        `title` VARCHAR(200) NOT NULL,
        `overview` TEXT NOT NULL,
        `release_date` DATE NOT NULL,
        `runtime_minutes` INT NOT NULL CHECK (runtime_minutes > 0),
        `poster_url` VARCHAR(50),
        `series_id` INT,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (series_id) REFERENCES `series` (id) ON DELETE SET NULL,
        INDEX idx_series_id (series_id),
        INDEX idx_release_date (release_date)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 映画ジャンル中間テーブル
CREATE TABLE
    `movie_genres` (
        `movie_id` INT NOT NULL,
        `genre_id` INT NOT NULL,
        PRIMARY KEY (movie_id, genre_id),
        FOREIGN KEY (movie_id) REFERENCES `movies` (id) ON DELETE CASCADE,
        FOREIGN KEY (genre_id) REFERENCES `genres` (id) ON DELETE CASCADE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- プラットフォームテーブル
CREATE TABLE
    `platforms` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `code` VARCHAR(50) NOT NULL UNIQUE,
        `name` VARCHAR(100) NOT NULL
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- ウォッチリストテーブル
CREATE TABLE
    `watchlist` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `user_id` INT NOT NULL,
        `movie_id` INT NOT NULL,
        `priority` TINYINT DEFAULT 1 CHECK (priority BETWEEN 1 AND 5) COMMENT '1が優先度高',
        `added_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        UNIQUE KEY unique_user_movie (user_id, movie_id),
        FOREIGN KEY (user_id) REFERENCES `users` (id) ON DELETE CASCADE,
        FOREIGN KEY (movie_id) REFERENCES `movies` (id) ON DELETE CASCADE,
        INDEX idx_added_at (added_at)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- レビューテーブル
CREATE TABLE
    `reviews` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `user_id` INT NOT NULL,
        `movie_id` INT NOT NULL,
        `rating` DECIMAL(2, 1) CHECK (rating BETWEEN 0.1 AND 5.0),
        `comment` TEXT,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        UNIQUE KEY unique_user_movie_review (user_id, movie_id),
        FOREIGN KEY (user_id) REFERENCES `users` (id) ON DELETE CASCADE,
        FOREIGN KEY (movie_id) REFERENCES `movies` (id) ON DELETE CASCADE,
        INDEX idx_created_at (created_at)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

-- 視聴履歴テーブル
CREATE TABLE
    `watch_history` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `user_id` INT NOT NULL,
        `movie_id` INT NOT NULL,
        `platform_id` INT NOT NULL,
        `watched_date` DATE,
        `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES `users` (id) ON DELETE CASCADE,
        FOREIGN KEY (movie_id) REFERENCES `movies` (id) ON DELETE CASCADE,
        FOREIGN KEY (platform_id) REFERENCES `platforms` (id) ON DELETE CASCADE,
        INDEX idx_user_watched (user_id, watched_date),
        INDEX idx_watched_date (watched_date)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;

USE db4test;

-- db4devのテーブル構造をコピー
CREATE TABLE
    `users` LIKE `db4dev`.`users`;

CREATE TABLE
    `refresh_tokens` LIKE `db4dev`.`refresh_tokens`;

CREATE TABLE
    `genres` LIKE `db4dev`.`genres`;

CREATE TABLE
    `series` LIKE `db4dev`.`series`;

CREATE TABLE
    `movies` LIKE `db4dev`.`movies`;

CREATE TABLE
    `movie_genres` LIKE `db4dev`.`movie_genres`;

CREATE TABLE
    `platforms` LIKE `db4dev`.`platforms`;

CREATE TABLE
    `watchlist` LIKE `db4dev`.`watchlist`;

CREATE TABLE
    `reviews` LIKE `db4dev`.`reviews`;

CREATE TABLE
    `watch_history` LIKE `db4dev`.`watch_history`;