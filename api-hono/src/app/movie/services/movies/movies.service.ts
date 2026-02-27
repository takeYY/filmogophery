import { Ok, ok } from "neverthrow";
import { getMoviesByGenre } from "../../repositories/movies/movies.repository";

export interface Genre {
  code: string;
  name: string;
}

export interface Movie {
  id: number;
  title: string;
  overview: string;
  releaseDate: string;
  runtimeMinute: number;
  posterUrl: string | null;
  tmdbId: number;
  genres: Genre[];
}

export async function getMovies(
  genre: string | undefined,
  limit: number,
  offset: number,
): Promise<Ok<Movie[], never>> {
  const result = await getMoviesByGenre(genre, limit, offset);
  if (result == undefined) {
    return ok([]);
  }

  const movies: Movie[] = result.map((movie) => {
    const codes = movie.genreCodes?.split(",") || [];
    const names = movie.genreNames?.split(",") || [];
    const genres = codes.map((code, i) => ({ code, name: names[i] }));

    return {
      id: movie.id,
      title: movie.title,
      overview: movie.overview,
      releaseDate: movie.releaseDate,
      runtimeMinute: movie.runtimeMinute,
      posterUrl: movie.posterUrl,
      tmdbId: movie.tmdbId,
      genres,
    };
  });

  return ok(movies);
}
