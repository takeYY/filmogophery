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

export type TmdbMovieResult = {
  id: number;
  title: string;
  overview: string;
  release_date: string;
  poster_path: string | null;
  genre_ids: number[];
};

export type TmdbSearchResponse = {
  page: number;
  results: TmdbMovieResult[];
  total_pages: number;
  total_results: number;
};

export type TmdbTrendingMovieResult = {
  id: number;
  title: string;
  overview: string;
  release_date: string;
  poster_path: string | null;
  genre_ids: number[];
};

export type TmdbTrendingResponse = {
  page: number;
  results: TmdbTrendingMovieResult[];
  total_pages: number;
  total_results: number;
};
