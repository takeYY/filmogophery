import { environment } from "@/core/environment";
import { err, ok } from "neverthrow";

type tmdbBelongsToCollection = {
  id: number;
  name: string;
  poster_path: string;
  backdrop_path: string;
};

type TmdbGenre = {
  id: number;
  name: string;
};

type tmdbProductionCompanies = {
  id: number;
  logo_path: string;
  name: string;
  origin_country: string;
};

type tmdbProductionCountries = {
  iso_3166_1: string;
  name: string;
};

type tmdbSpokenLanguages = {
  english_name: string;
  iso_639_1: string;
  name: string;
};

export type TmdbMovieDetailResponse = {
  id: number;
  adult: boolean;
  backdrop_path: string;
  original_language: string;
  original_title: string;
  overview: string;
  popularity: number;
  release_date: string;
  title: string;
  video: boolean;
  vote_average: number;
  vote_count: number;
  belongs_to_collection: tmdbBelongsToCollection;
  budget: number;
  genres: TmdbGenre[];
  homepage: string;
  imdb_id: string;
  origin_country: string[];
  poster_path: string;
  production_companies: tmdbProductionCompanies[];
  production_countries: tmdbProductionCountries[];
  revenue: number;
  runtime: number;
  spoken_languages: tmdbSpokenLanguages[];
  status: string;
  tagline: string;
};

const BASE_URL = `https://api.themoviedb.org/3`;

export async function getTmdbMovieDetailById(tmdbId: number) {
  const url = BASE_URL + `/movie/${tmdbId}?language=ja-JP`;

  try {
    const response = await fetch(url, {
      method: "GET",
      headers: {
        Accept: "application/json",
        Authorization: `Bearer ${environment.TMDB.ACCESS_TOKEN}`,
      },
    });
    return ok(response);
  } catch (e) {
    return err(e as TypeError | Error);
  }
}
