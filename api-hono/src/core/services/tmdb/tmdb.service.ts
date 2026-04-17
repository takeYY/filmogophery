import { environment } from "@/core/environment";
import { TmdbSearchResponse } from "@/core/types/tmdb";
import { err, ok } from "neverthrow";

const BASE_URL = `https://api.themoviedb.org/3`;

export async function getMoviesByTitle(title: string, offset: number) {
  const solid = 20;
  const page = offset / solid + 1;
  const url =
    BASE_URL + `/search/movie?language=ja-JP&query=${title}&page=${page}`;

  try {
    const response = await fetch(url, {
      method: "GET",
      headers: {
        Accept: "application/json",
        Authorization: `Bearer ${environment.TMDB.ACCESS_TOKEN}`,
      },
    });
    return ok((await response.json()) as TmdbSearchResponse);
  } catch (e) {
    return err(e as TypeError | Error);
  }
}

export async function getTmdbMovieDetailById(tmdbId: number) {
  const url = BASE_URL + `/movie/${tmdbId}?language=ja-JP`;

  try {
    const response = await fetch(url, {
      method: "GET",
      headers: {
        Accept: "application/json",
        Authorization: `Bearer ${environment.TMDB.ACCESS_TOKEN}`,
      },
    });
    return ok(response);
  } catch (e) {
    return err(e as TypeError | Error);
  }
}
