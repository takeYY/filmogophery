import { mysqlTable, mysqlSchema, AnyMySqlColumn, primaryKey, unique, int, varchar, index, foreignKey, check, text, date, timestamp, decimal, tinyint } from "drizzle-orm/mysql-core"
import { sql } from "drizzle-orm"

export const genres = mysqlTable("genres", {
	id: int().autoincrement().notNull(),
	code: varchar({ length: 50 }).notNull(),
	name: varchar({ length: 100 }).notNull(),
},
(table) => [
	primaryKey({ columns: [table.id], name: "genres_id"}),
	unique("code").on(table.code),
]);

export const movieGenres = mysqlTable("movie_genres", {
	movieId: int("movie_id").notNull().references(() => movies.id, { onDelete: "cascade" } ),
	genreId: int("genre_id").notNull().references(() => genres.id, { onDelete: "cascade" } ),
},
(table) => [
	index("genre_id").on(table.genreId),
	primaryKey({ columns: [table.movieId, table.genreId], name: "movie_genres_movie_id_genre_id"}),
]);

export const movies = mysqlTable("movies", {
	id: int().autoincrement().notNull(),
	tmdbId: int("tmdb_id").notNull(),
	title: varchar({ length: 200 }).notNull(),
	overview: text().notNull(),
	// you can use { mode: 'date' }, if you want to have Date as type for this column
	releaseDate: date("release_date", { mode: 'string' }).notNull(),
	runtimeMinutes: int("runtime_minutes").notNull(),
	posterUrl: varchar("poster_url", { length: 50 }),
	seriesId: int("series_id").references(() => series.id, { onDelete: "set null" } ),
	createdAt: timestamp("created_at", { mode: 'string' }).defaultNow(),
	updatedAt: timestamp("updated_at", { mode: 'string' }).defaultNow().onUpdateNow(),
},
(table) => [
	index("idx_release_date").on(table.releaseDate),
	index("idx_series_id").on(table.seriesId),
	primaryKey({ columns: [table.id], name: "movies_id"}),
	unique("tmdb_id").on(table.tmdbId),
	check("movies_chk_1", sql`(\`runtime_minutes\` > 0)`),
]);

export const platforms = mysqlTable("platforms", {
	id: int().autoincrement().notNull(),
	code: varchar({ length: 50 }).notNull(),
	name: varchar({ length: 100 }).notNull(),
},
(table) => [
	primaryKey({ columns: [table.id], name: "platforms_id"}),
	unique("code").on(table.code),
]);

export const pointHistory = mysqlTable("point_history", {
	id: int().autoincrement().notNull(),
	userId: int("user_id").notNull().references(() => users.id, { onDelete: "cascade" } ),
	points: int().notNull(),
	action: varchar({ length: 50 }).notNull(),
	referenceId: int("reference_id").notNull(),
	createdAt: timestamp("created_at", { mode: 'string' }).defaultNow(),
},
(table) => [
	index("idx_created_at").on(table.createdAt),
	index("idx_user_id").on(table.userId),
	primaryKey({ columns: [table.id], name: "point_history_id"}),
]);

export const refreshTokens = mysqlTable("refresh_tokens", {
	id: int().autoincrement().notNull(),
	userId: int("user_id").notNull().references(() => users.id, { onDelete: "cascade" } ),
	tokenHash: varchar("token_hash", { length: 255 }).notNull(),
	expiresAt: timestamp("expires_at", { mode: 'string' }).notNull(),
	createdAt: timestamp("created_at", { mode: 'string' }).defaultNow(),
	revokedAt: timestamp("revoked_at", { mode: 'string' }),
},
(table) => [
	index("idx_expires_at").on(table.expiresAt),
	index("idx_user_id").on(table.userId),
	primaryKey({ columns: [table.id], name: "refresh_tokens_id"}),
	unique("token_hash").on(table.tokenHash),
]);

export const reviews = mysqlTable("reviews", {
	id: int().autoincrement().notNull(),
	userId: int("user_id").notNull().references(() => users.id, { onDelete: "cascade" } ),
	movieId: int("movie_id").notNull().references(() => movies.id, { onDelete: "cascade" } ),
	rating: decimal({ precision: 2, scale: 1 }),
	comment: text(),
	createdAt: timestamp("created_at", { mode: 'string' }).defaultNow(),
	updatedAt: timestamp("updated_at", { mode: 'string' }).defaultNow().onUpdateNow(),
},
(table) => [
	index("idx_created_at").on(table.createdAt),
	index("movie_id").on(table.movieId),
	primaryKey({ columns: [table.id], name: "reviews_id"}),
	unique("unique_user_movie_review").on(table.userId, table.movieId),
	check("reviews_chk_1", sql`(\`rating\` between 0.1 and 5.0)`),
]);

export const series = mysqlTable("series", {
	id: int().autoincrement().notNull(),
	name: varchar({ length: 200 }).notNull(),
	posterUrl: varchar("poster_url", { length: 50 }),
	createdAt: timestamp("created_at", { mode: 'string' }).defaultNow(),
	updatedAt: timestamp("updated_at", { mode: 'string' }).defaultNow().onUpdateNow(),
},
(table) => [
	primaryKey({ columns: [table.id], name: "series_id"}),
	unique("name").on(table.name),
]);

export const userPoints = mysqlTable("user_points", {
	id: int().autoincrement().notNull(),
	userId: int("user_id").notNull().references(() => users.id, { onDelete: "cascade" } ),
	totalPoints: int("total_points").default(0).notNull(),
	level: tinyint().default(1).notNull(),
	createdAt: timestamp("created_at", { mode: 'string' }).defaultNow(),
	updatedAt: timestamp("updated_at", { mode: 'string' }).defaultNow().onUpdateNow(),
},
(table) => [
	primaryKey({ columns: [table.id], name: "user_points_id"}),
	unique("user_id").on(table.userId),
]);

export const users = mysqlTable("users", {
	id: int().autoincrement().notNull(),
	username: varchar({ length: 50 }).notNull(),
	email: varchar({ length: 100 }).notNull(),
	passwordHash: varchar("password_hash", { length: 255 }).notNull(),
	isActive: tinyint("is_active").default(1),
	lastLoginAt: timestamp("last_login_at", { mode: 'string' }),
	createdAt: timestamp("created_at", { mode: 'string' }).defaultNow(),
	updatedAt: timestamp("updated_at", { mode: 'string' }).defaultNow().onUpdateNow(),
},
(table) => [
	primaryKey({ columns: [table.id], name: "users_id"}),
	unique("email").on(table.email),
	unique("username").on(table.username),
]);

export const watchHistory = mysqlTable("watch_history", {
	id: int().autoincrement().notNull(),
	userId: int("user_id").notNull().references(() => users.id, { onDelete: "cascade" } ),
	movieId: int("movie_id").notNull().references(() => movies.id, { onDelete: "cascade" } ),
	platformId: int("platform_id").notNull().references(() => platforms.id, { onDelete: "cascade" } ),
	// you can use { mode: 'date' }, if you want to have Date as type for this column
	watchedDate: date("watched_date", { mode: 'string' }),
	createdAt: timestamp("created_at", { mode: 'string' }).defaultNow(),
},
(table) => [
	index("idx_user_watched").on(table.userId, table.watchedDate),
	index("idx_watched_date").on(table.watchedDate),
	index("movie_id").on(table.movieId),
	index("platform_id").on(table.platformId),
	primaryKey({ columns: [table.id], name: "watch_history_id"}),
]);

export const watchlist = mysqlTable("watchlist", {
	id: int().autoincrement().notNull(),
	userId: int("user_id").notNull().references(() => users.id, { onDelete: "cascade" } ),
	movieId: int("movie_id").notNull().references(() => movies.id, { onDelete: "cascade" } ),
	priority: tinyint().default(1),
	addedAt: timestamp("added_at", { mode: 'string' }).defaultNow(),
	updatedAt: timestamp("updated_at", { mode: 'string' }).defaultNow().onUpdateNow(),
},
(table) => [
	index("idx_added_at").on(table.addedAt),
	index("movie_id").on(table.movieId),
	primaryKey({ columns: [table.id], name: "watchlist_id"}),
	unique("unique_user_movie").on(table.userId, table.movieId),
	check("watchlist_chk_1", sql`(\`priority\` between 1 and 5)`),
]);
