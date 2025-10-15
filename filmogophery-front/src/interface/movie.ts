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
