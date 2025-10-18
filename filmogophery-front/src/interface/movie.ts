export type Genre = {
  code: string;
  name: string;
};

type Poster = {
  id: number;
  url: string;
};

export type SeriesNeo = {
  id: number;
  name: string;
  poster_id: number;
  poster: Poster | null;
};

export type MovieImpression = {
  id: number;
  movie_id: number;
  status: boolean;
  rating: number;
  note: string;
  movie: Movie | null;
  watch_records: null;
};

// deprecated
export type Movie = {
  id: number;
  title: string;
  overview: string;
  releaseDate: string;
  runTime: number;
  posterURL: string | null;
  tmdbID: number;
  genres: Genre[];
};

export type MovieNeo = {
  id: number;
  title: string;
  overview: string;
  releaseDate: string;
  runtimeMinutes: number;
  posterURL: string | null;
  tmdbID: number;
  genres: Genre[];
};

export type MovieDetailNeo = {
  id: number;
  title: string;
  overview: string;
  releaseDate: string;
  runtimeMinutes: number;
  posterURL: string | null;
  tmdbID: number;
  voteAverage: number; // Min:0.0, Max:5.0
  voteCount: number;
  genres: Genre[];
  series: SeriesNeo2 | null;
  review: Review | null;
};

type SeriesNeo2 = {
  name: string;
  posterURL: string | null;
};

export type Review = {
  id: number;
  rating: number | null;
  comment: string | null;
  createdAt: string;
  updatedAt: string;
};

export type WatchHistory = {
  id: number;
  platform: Platform;
  watchedAt: string;
};

type Platform = {
  code: string;
  name: string;
};

// deprecated
type Series = {
  name: string;
  posterURL: string;
};

export type Impression = {
  id: number;
  status: string;
  rating: number | null; // Min:0.0, Max:5.0
  note: string;
  records: WatchRecord[];
};

export type ImpressionResult = {
  rating: number | null; // Min:0.0, Max:5.0
  note: string;
};

// deprecated
export type WatchRecord = {
  watchDate: string;
  watchMedia: string;
};

// deprecated
export type MovieDetail = {
  id: number;
  title: string;
  overview: string | null;
  releaseDate: string;
  runTime: number;
  posterURL: string;
  tmdbID: number;
  voteAverage: number; // Min:0.0, Max:5.0
  voteCount: number;
  genres: Genre[];
  series: Series | null;
  impression: Impression | null;
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
