USE db4dev;

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- Create `users`
INSERT INTO
    `users` (`id`, `username`, `email`)
VALUES (1, 'システム管理者', 'system@filmogophery.jp');

-- Create `genres`
INSERT INTO
    `genres` (`id`, `code`, `name`)
VALUES
    (12, 'adventure', 'アドベンチャー'),
    (14, 'fantasy', 'ファンタジー'),
    (16, 'animation', 'アニメーション'),
    (18, 'drama', 'ドラマ'),
    (27, 'horror', 'ホラー'),
    (28, 'action', 'アクション'),
    (35, 'comedy', 'コメディ'),
    (36, 'history', 'ヒストリー'),
    (37, 'western', '西部劇'),
    (53, 'thriller', 'スリラー'),
    (80, 'crime', 'クライム'),
    (99, 'documentary', 'ドキュメンタリー'),
    (878, 'sf', 'SF'),
    (9648, 'mystery', 'ミステリー'),
    (10402, 'music', 'ミュージック'),
    (10749, 'romance', 'ロマンス'),
    (10751, 'family', 'ファミリー'),
    (10752, 'war', '戦争'),
    (10770, 'tv', 'TV');

-- Create `platforms`
INSERT INTO
    `platforms` (`id`, `code`, `name`)
VALUES
    (1, 'primeVideo', 'Prime Video'),
    (2, 'netflix', 'Netflix'),
    (3, 'uNext', 'U-NEXT'),
    (4, 'disneyPlus', 'Disney+'),
    (5, 'youtube', 'YouTube'),
    (6, 'appleTv', 'Apple TV+'),
    (7, 'hulu', 'Hulu'),
    (8, 'dAnime', 'dアニメ'),
    (9, 'telasa', 'TELASA'),
    (10, 'cinema', '映画館'),
    (99, 'unknown', '不明');

-- Create `series`
INSERT INTO
    `series` (`id`, `name`, `poster_url`)
VALUES
    (
        1,
        'ターミネーターシリーズ',
        '/pF5GIijY2fyZcByqNDzhS8v4h1x.jpg'
    ),
    (2, 'バック・トゥ・ザ・フューチャー シリーズ', NULL);

-- Create `movies`
INSERT INTO
    `movies` (
        `id`,
        `title`,
        `overview`,
        `release_date`,
        `runtime_minutes`,
        `poster_url`,
        `series_id`,
        `tmdb_id`
    )
VALUES
    (
        218, -- id
        'ターミネーター', -- title
        'アメリカのとある街、深夜突如奇怪な放電と共に屈強な肉体をもった男が現れる。同じく...', -- overview
        '1985-05-04', -- release_date
        108, -- runtime_minutes
        '/iK2YBfD7DdaNXZALQhhT9uTN9Rc.jpg', -- poster_url
        1, -- series_id
        218 -- tmdb_id
    ),
    (
        280, -- id
        'ターミネーター2', -- title
        '未来からの抹殺兵器ターミネーターを破壊し、近未来で恐ろしい戦争が起こる事を知って...', -- overview
        '1991-08-24', -- release_date
        137, -- runtime_minutes
        '/oCwo7ALD3LftLqy0Oj6U669u4fU.jpg', -- poster_url
        1, -- series_id
        280 -- tmdb_id
    ),
    (
        296, -- id
        'ターミネーター3', -- title
        'スカイネットを破壊、予見されていた最終戦争の日も過ぎたが、ジョンの心には母親から...', -- overview
        '2003-07-12', -- release_date
        110, -- runtime_minutes
        NULL, -- poster_url
        1, -- series_id
        296 -- tmdb_id
    ),
    (
        534, -- id
        'ターミネーター4', -- title
        '“審判の日”から10年後の2018年。人類軍の指導者となり、機械軍と戦うことを幼...', -- overview
        '2009-06-05', -- release_date
        114, -- runtime_minutes
        '/7fmBDyHsLTyhfRSZk19wrY23zDg.jpg', -- poster_url
        1, -- series_id
        534 -- tmdb_id
    ),
    (
        105, -- id
        'バック・トゥ・ザ・フューチャー', -- title
        'スティーブン・スピルバーグとロバート・ゼメキスが贈るSFアドベンチャーシリーズ第...', -- overview
        '1985-12-07', -- release_date
        116, -- runtime_minutes
        NULL, -- poster_url
        2, -- series_id
        105 -- tmdb_id
    ),
    (
        19995, -- id
        'アバター', -- title
        '西暦2154年。人類は惑星ポリフェマスの最大衛星パンドラに鉱物採掘基地を開いてい...', -- overview
        '2010-08-26', -- release_date
        162, -- runtime_minutes
        NULL, -- poster_url
        NULL, -- series_id
        19995 -- tmdb_id
    );

-- Create `movie_genres`
INSERT INTO
    `movie_genres` (`movie_id`, `genre_id`)
VALUES
    (218, 28),
    (218, 27),
    (218, 878),
    (534, 28),
    (534, 10752),
    (534, 27),
    (534, 878),
    (534, 53),
    (105, 28),
    (105, 10402),
    (105, 35),
    (105, 18),
    (105, 10751),
    (105, 10749),
    (105, 878);

-- Create `reviews`
INSERT INTO
    `reviews` (`id`, `user_id`, `movie_id`, `rating`, `comment`)
VALUES
    (
        1,
        1,
        218,
        4.3,
        'ターミネーターの元祖という感じで、恐ろしさと希望が織り成す圧巻の作品。今観るとVFXのぎこちなさが目立つが、それが逆に怖さを演出している。'
    ),
    (2, 1, 296, NULL, NULL);

-- Create `watch_history`
INSERT INTO
    `watch_history` (
        `id`,
        `review_id`,
        `platform_id`,
        `watched_date`
    )
VALUES
    (1, 1, 99, '1895-12-28'),
    (2, 1, 1, '2022-10-24'),
    (3, 1, 3, '2024-08-01');


USE db4test;

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- Create `users`
INSERT INTO
    `users` (`id`, `username`, `email`)
VALUES (1, 'システム管理者', 'system@filmogophery.jp');

-- マスタデータをdb4devからコピー
INSERT INTO db4test.genres SELECT * FROM db4dev.genres;
INSERT INTO db4test.platforms SELECT * FROM db4dev.platforms;
