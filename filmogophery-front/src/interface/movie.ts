export type Genre = {
  id: number;
  code: string;
  name: string;
  movies: null;
};

type Poster = {
  id: number;
  url: string;
};

type SeriesNeo = {
  id: number;
  name: string;
  poster_id: number;
  poster: Poster;
};

type MovieImpression = {
  id: number;
  movie_id: number;
  status: boolean;
  rating: number;
  note: string;
  movie: Movie;
  watch_records: null;
};

export type Movie = {
  id: number;
  title: string;
  overview: string;
  release_date: string;
  run_time: number;
  poster_url: string;
  series_id: number;
  tmdb_id: number;
  genres: Genre[];
  series: SeriesNeo;
  movie_impression: MovieImpression;
};

type Series = {
  name: string;
  posterURL: string;
};

type Impression = {
  id: number;
  status: string;
  rating: number | null;
  note: string | null;
};

export type WatchRecord = {
  watchDate: string;
  watchMedia: string;
};

export type MovieDetail = {
  id: number;
  title: string;
  overview: string | null;
  releaseDate: string;
  runTime: number;
  genres: string[];
  posterURL: string;
  voteAverage: number;
  voteCount: number;
  series: Series | null;
  impression: Impression;
  watchRecords: WatchRecord[];
};

export type WatchMedia = {
  id: number;
  code: string;
  name: string;
};

export type SearchMovie = {
  tmdbID: number;
  title: string;
  overview: string;
  popularity: number;
  posterURL: string;
  releaseDate: string;
  voteAverage: number;
  voteCount: number;
  genres: string[];
};
