import { Movie } from "./movie";

export type Watchlist = {
  id: number;
  priority: number;
  addedAt: string;
  movie: Movie;
};
