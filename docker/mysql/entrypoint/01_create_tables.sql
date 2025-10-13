USE db4dev;

-- ユーザーテーブル
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- ジャンルテーブル
CREATE TABLE genres (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL
);

-- シリーズテーブル
CREATE TABLE series (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    poster_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 映画テーブル
CREATE TABLE movies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    overview TEXT NOT NULL,
    release_date DATE NOT NULL,
    runtime_minutes INT NOT NULL,
    poster_url VARCHAR(255),
    series_id INT,
    tmdb_id INT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE SET NULL
);

-- 映画ジャンル中間テーブル
CREATE TABLE movie_genres (
    movie_id INT NOT NULL,
    genre_id INT NOT NULL,
    PRIMARY KEY (movie_id, genre_id),
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE
);

-- プラットフォームテーブル
CREATE TABLE platforms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL
);

-- ウォッチリストテーブル
CREATE TABLE watchlist (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    movie_id INT NOT NULL,
    priority TINYINT DEFAULT 1,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_movie (user_id, movie_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE
);

-- レビューテーブル
CREATE TABLE reviews (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    movie_id INT NOT NULL,
    rating DECIMAL(2,1),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_movie_review (user_id, movie_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE
);

-- 視聴履歴テーブル
CREATE TABLE watch_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    review_id INT NOT NULL,
    platform_id INT NOT NULL,
    watched_date DATE NOT NULL DEFAULT "1895-12-28",
    FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
    FOREIGN KEY (platform_id) REFERENCES platforms(id) ON DELETE RESTRICT
);


USE db4test;

-- db4devのテーブル構造をコピー
CREATE TABLE users LIKE db4dev.users;
CREATE TABLE genres LIKE db4dev.genres;
CREATE TABLE series LIKE db4dev.series;
CREATE TABLE movies LIKE db4dev.movies;
CREATE TABLE movie_genres LIKE db4dev.movie_genres;
CREATE TABLE platforms LIKE db4dev.platforms;
CREATE TABLE watchlist LIKE db4dev.watchlist;
CREATE TABLE reviews LIKE db4dev.reviews;
CREATE TABLE watch_history LIKE db4dev.watch_history;
