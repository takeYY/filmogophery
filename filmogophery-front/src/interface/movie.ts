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
  poster_id: number;
  series_id: number;
  tmdb_id: number;
  genres: Genre[];
  poster: Poster | null;
  series: SeriesNeo;
  movie_impression: MovieImpression;
};

type Series = {
  name: string;
  posterURL: string;
};

type Impression = {
  status: string;
  rating: number;
  note: string | null;
};

export type WatchRecord = {
  watch_media: string;
  watch_date: Date;
};

export type MovieDetail = {
  id: number;
  title: string;
  overview: string | null;
  release_date: Date;
  run_time: number;
  genres: string[];
  posterURL: string;
  vote_average: number;
  vote_count: number;
  series: Series | null;
  impression: Impression;
  watchRecords: WatchRecord[];
};

export type WatchMedia = {
  id: number;
  code: string;
  name: string;
};
