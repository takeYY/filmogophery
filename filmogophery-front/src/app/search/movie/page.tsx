"use client";

import React, { useEffect, useState } from "react";
import { SearchMovie } from "@/interface/movie";
import { useSearchParams } from "next/navigation";
import Image from "next/image";
import { posterUrlPrefix } from "@/constants/poster";
import StarRating from "@/app/components/Rating";
import { APIBaseURL } from "@/constants/api";
import Link from "next/link";

export default function Page() {
  const searchParams = useSearchParams();
  const query = searchParams.get("query");

  const [movies, setMovies] = useState<SearchMovie[]>([]);

  useEffect(() => {
    const fetchMovie = async () => {
      try {
        // TODO: 何故か Dynamic Routing が効かないので、後で直すこと!!
        //const response = await fetch(`/api/search/movie?query=${query}`, {
        const response = await fetch(
          `${APIBaseURL}/tmdb/search/movies?query=${query}`,
          {
            method: "GET",
          }
        );
        const movies: SearchMovie[] = await response.json();
        console.log("moviesのデータ取得: 完了");
        console.log("%o", movies);

        return setMovies(movies);
      } catch {
        console.log("moviesのデータ取得: エラー。空配列で定義します");
        return setMovies([]);
      }
    };
    fetchMovie();
  }, [query]);

  async function addWatchList(movie: SearchMovie) {
    console.log(`movie is %o`, movie);
    try {
      const response = await fetch(`/api/movie`, {
        method: "POST",
        body: JSON.stringify({
          tmdbID: movie.tmdbID,
          status: false,
        }),
      });
      const status = response.status;
      if (status === 201) {
        console.log("success to add watch list");
      } else {
        console.log("failed to add watch list");
      }
    } catch {
      console.log("failed to add watch list");
    }
  }

  return (
    <div className="container-fluid pb-4">
      <h3 className="text-center mb-4">Search Movies</h3>

      <div className="row row-cols-md-3 g-4">
        {movies.map((movie: SearchMovie, index: number) => {
          return (
            <div className="col">
              <div className="card mb-2 bg-dark border-info">
                <div className="row g-0">
                  <div className="col-md-4">
                    {/* ポスター */}
                    <Image
                      src={
                        posterUrlPrefix +
                        (movie.posterURL
                          ? movie.posterURL
                          : "/Agz71U0wcesx87micVn731Z1vPu.jpg")
                      }
                      className="card-img-top"
                      alt="..."
                      width={250}
                      height={250}
                    />

                    {/* 一般の評価 */}
                    <div className="justify-content-center">
                      <StarRating
                        rating={movie.voteAverage / 2}
                        size={20}
                        starColor={"#0dcaf0"}
                        sumReview={movie.voteCount.toString()}
                      />
                    </div>
                  </div>
                  <div className="col-md-8">
                    <div className="card-body text-light">
                      {/* タイトル */}
                      <h5 className="card-title">{movie.title}</h5>
                      {/* ジャンル */}
                      {movie.genres.length !== 0 && (
                        <div className="card-text d-grid gap-2 d-md-block">
                          {movie.genres.map((g: string, i: number) => {
                            return (
                              <button
                                key={i}
                                type="button"
                                className="btn btn-outline-info btn-sm"
                              >
                                {g}
                              </button>
                            );
                          })}
                        </div>
                      )}
                      {/* 公開日 */}
                      <p className="card-text">公開日：{movie.releaseDate}</p>
                      {/* 概要 */}
                      <p className="card-text">
                        {movie.overview.length > 40
                          ? movie.overview.substring(0, 37) + "..."
                          : movie.overview}
                      </p>

                      <div className="border-top border-success">
                        <div className="row mt-2">
                          <div className="col-md-6 text-center">
                            <button
                              className="btn btn-outline-warning"
                              onClick={() => addWatchList(movie)}
                            >
                              Watch List
                            </button>
                          </div>
                          <div className="col-md-6 text-center">
                            <Link className="btn btn-outline-success" href="">
                              Note
                            </Link>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
