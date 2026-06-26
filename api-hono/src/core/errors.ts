export class MovieIsNotFound extends Error {
  constructor(id?: number) {
    super(
      id !== undefined ? `movie(id=${id}) is not found` : "movie is not found",
    );
    this.name = "MovieIsNotFound";
  }
}

export class TmdbMovieIsNotFound extends Error {
  constructor(tmdbId?: number) {
    super(
      tmdbId !== undefined
        ? `tmdb movie(id=${tmdbId}) is not found`
        : "tmdb movie is not found",
    );
    this.name = "TmdbMovieIsNotFound";
  }
}
