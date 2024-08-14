export type Movie = {
  id: number;
  title: string;
  overview: string | null;
  release_date: string;
  run_time: number;
  genres: string[];
  posterURL: string;
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
