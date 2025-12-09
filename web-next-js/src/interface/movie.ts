import { Genre } from "./genre";
import { Review } from "./review";

export type Movie = {
  id: number;
  title: string;
  overview: string;
  releaseDate: string;
  runtimeMinutes: number;
  posterURL: string | null;
  tmdbID: number;
  genres: Genre[];
};

// Movie を継承
export type MovieDetail = Movie & {
  voteAverage: number; // Min:0.0, Max:5.0
  voteCount: number;
  series: Series | null;
  review: Review | null;
};

type Series = {
  name: string;
  posterURL: string | null;
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
