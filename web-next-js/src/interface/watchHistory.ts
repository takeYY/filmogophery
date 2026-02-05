import { Movie } from "./movie";
import { Platform } from "./platform";

export type WatchHistory = {
  id: number;
  platform: Platform;
  watchedAt: string;
};

export type MyWatchHistory = {
  id: number;
  platform: Platform;
  watchedAt: string;
  movie: Movie;
};
