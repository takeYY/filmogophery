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

export interface Review {
  id: number;
  createdAt: string | null;
  updatedAt: string | null;
  rating: string | null;
  comment: string | null;
}

export interface Series {
  name: string | null;
  posterUrl: string | null;
}

export interface MovieDetail extends Movie {
  voteAverage: number;
  voteCount: number;
  series: Series | null;
  review: Review | null;
}
