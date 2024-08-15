-- Create `genre`
INSERT INTO
    `genre` (`id`, `code`, `name`)
VALUES
    (1, 'action', 'アクション'),
    (2, 'adventure', 'アドベンチャー'),
    (3, 'animation', 'アニメーション'),
    (4, 'comedy', 'コメディ'),
    (5, 'crime', 'クライム'),
    (6, 'documentary', 'ドキュメンタリー'),
    (7, 'drama', 'ドラマ'),
    (8, 'family', 'ファミリー'),
    (9, 'fantasy', 'ファンタジー'),
    (10, 'history', 'ヒストリー'),
    (11, 'horror', 'ホラー'),
    (12, 'musical', 'ミュージカル'),
    (13, 'mystery', 'ミステリー'),
    (14, 'romance', 'ロマンス'),
    (15, 'sf', 'SF'),
    (16, 'tv', 'TV'),
    (17, 'thriller', 'スリラー'),
    (18, 'war', '戦争'),
    (19, 'western', '西部劇');

-- Create `watch_media`
INSERT INTO
    `watch_media` (`id`, `code`, `name`)
VALUES
    (1, 'prime_video', 'Prime Video'),
    (2, 'netflix', 'Netflix'),
    (3, 'u_next', 'U-NEXT'),
    (4, 'disney_plus', 'Disney+'),
    (5, 'youtube', 'YouTube'),
    (6, 'apple_tv', 'Apple TV+'),
    (7, 'hulu', 'Hulu'),
    (8, 'd_anime', 'dアニメ'),
    (9, 'telasa', 'TELASA'),
    (10, 'cinema', '映画館'),
    (99, 'unknown', '不明');

-- Create `poster`
INSERT INTO
    `poster` (`id`, `url`)
VALUES
    (1, '/iK2YBfD7DdaNXZALQhhT9uTN9Rc.jpg'), -- ターミネーター
    (2, '/oCwo7ALD3LftLqy0Oj6U669u4fU.jpg'), -- ターミネーター2
    (3, '/7fmBDyHsLTyhfRSZk19wrY23zDg.jpg'), -- ターミネーター4
    (4, '/pF5GIijY2fyZcByqNDzhS8v4h1x.jpg') -- ターミネーターシリーズ
;

-- Create `movie_series`
INSERT INTO
    `movie_series` (`id`, `name`, `poster_id`)
VALUES
    (1, 'ターミネーターシリーズ', 4),
    (2, 'バック・トゥ・ザ・フューチャー シリーズ', NULL);

-- Create `movie`
INSERT INTO
    `movie` (
        `id`,
        `title`,
        `overview`,
        `release_date`,
        `run_time`,
        `poster_id`,
        `series_id`,
        `tmdb_id`
    )
VALUES
    (
        1, -- id
        'ターミネーター', -- title
        'アメリカのとある街、深夜突如奇怪な放電と共に屈強な肉体をもった男が現れる。同じく...', -- overview
        '1985-05-04', -- release_date
        108, -- run_time
        1, -- poster_id
        1, -- series_id
        218 -- tmdb_id
    ),
    (
        2, -- id
        'ターミネーター2', -- title
        '未来からの抹殺兵器ターミネーターを破壊し、近未来で恐ろしい戦争が起こる事を知って...', -- overview
        '1991-08-24', -- release_date
        137, -- run_time
        2, -- poster_id
        1, -- series_id
        280 -- tmdb_id
    ),
    (
        3, -- id
        'ターミネーター3', -- title
        'スカイネットを破壊、予見されていた最終戦争の日も過ぎたが、ジョンの心には母親から...', -- overview
        '2003-07-12', -- release_date
        110, -- run_time
        NULL, -- poster_id
        1, -- series_id
        296 -- tmdb_id
    ),
    (
        4, -- id
        'ターミネーター4', -- title
        '“審判の日”から10年後の2018年。人類軍の指導者となり、機械軍と戦うことを幼...', -- overview
        '2009-06-05', -- release_date
        114, -- run_time
        3, -- poster_id
        1, -- series_id
        534 -- tmdb_id
    ),
    (
        5, -- id
        'バック・トゥ・ザ・フューチャー', -- title
        'スティーブン・スピルバーグとロバート・ゼメキスが贈るSFアドベンチャーシリーズ第...', -- overview
        '1985-12-07', -- release_date
        116, -- run_time
        NULL, -- poster_id
        2, -- series_id
        105 -- tmdb_id
    ),
    (
        6, -- id
        'アバター', -- title
        '西暦2154年。人類は惑星ポリフェマスの最大衛星パンドラに鉱物採掘基地を開いてい...', -- overview
        '2010-08-26', -- release_date
        162, -- run_time
        NULL, -- poster_id
        NULL, -- series_id
        19995 -- tmdb_id
    );

-- Create `movie_genres`
INSERT INTO
    `movie_genres` (`movie_id`, `genre_id`)
VALUES
    (1, 1),
    (1, 11),
    (1, 15),
    (4, 1),
    (4, 7),
    (4, 11),
    (4, 15),
    (4, 18),
    (5, 1),
    (5, 2),
    (5, 4),
    (5, 7),
    (5, 8),
    (5, 14),
    (5, 15);

-- Create `movie_impression`
INSERT INTO
    `movie_impression` (`id`, `movie_id`, `status`, `rating`, `note`)
VALUES
    (
        1,
        1,
        1,
        4.3,
        'ターミネーターの元祖という感じで、恐ろしさと希望が織り成す圧巻の作品。今観るとCGのぎこちなさが目立つが、それが逆に怖さを演出している。'
    ),
    (2, 3, 0, NULL, NULL);

-- Create `movie_watch_record`
INSERT INTO
    `movie_watch_record` (
        `id`,
        `movie_impression_id`,
        `watch_media_id`,
        `watch_date`
    )
VALUES
    (1, 1, 99, '2016-12-25'),
    (2, 1, 1, '2022-10-24'),
    (3, 1, 3, '2024-08-01');