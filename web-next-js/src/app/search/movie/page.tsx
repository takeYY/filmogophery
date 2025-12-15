"use client";

import React, { useEffect, useState, useCallback } from "react";
import { Movie, Genre } from "@/interface/index";
import { useSearchParams } from "next/navigation";
import Image from "next/image";
import { posterUrlPrefix } from "@/constants/poster";
import Link from "next/link";

export default function SearchMovies() {
  const searchParams = useSearchParams();
  const query = searchParams.get("query");

  const [movies, setMovies] = useState<Movie[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchMovie = async () => {
      if (!query) {
        setLoading(false);
        return;
      }

      setLoading(true);
      try {
        const response = await fetch(`/api/search/movie?query=${query}`, {
          method: "GET",
        });
        const movies: Movie[] = await response.json();
        console.log("moviesのデータ取得: 完了");
        console.log("%o", movies);

        setMovies(movies);
      } catch {
        console.log("moviesのデータ取得: エラー。空配列で定義します");
        setMovies([]);
      } finally {
        setLoading(false);
      }
    };
    fetchMovie();
  }, [query]);

  const addWatchList = useCallback(async (movie: Movie) => {
    console.log(`movie is %o`, movie);
    console.log("successfully added a movie to watchlist");
    try {
    } catch {
      console.log("failed to add watch list");
    }
  }, []);

  return (
    <div className="container-fluid pb-4">
      <h3 className="text-center mb-4">Search Movies</h3>

      {loading ? (
        <div className="text-center">
          <div className="spinner-border text-info" role="status">
            <span className="visually-hidden">Loading...</span>
          </div>
        </div>
      ) : movies.length === 0 ? (
        <div className="text-center text-light">
          <p>検索結果が見つかりませんでした。</p>
        </div>
      ) : (
        <div className="row row-cols-md-3 g-4">
          {movies.map((movie: Movie) => {
            return (
              <div className="col" key={movie.id}>
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
                        className="card-img-top w-100 h-auto"
                        alt="..."
                        width={250}
                        height={375}
                        style={{ objectFit: "cover" }}
                      />
                    </div>
                    <div className="col-md-8">
                      <div className="card-body text-light">
                        {/* タイトル */}
                        <h5 className="card-title">{movie.title}</h5>
                        {/* ジャンル */}
                        {movie.genres.length !== 0 && (
                          <div className="card-text d-grid gap-2 d-md-block">
                            {movie.genres.map((g: Genre, i: number) => {
                              return (
                                <button
                                  key={g.code}
                                  type="button"
                                  className="btn btn-outline-info btn-sm"
                                >
                                  {g.name}
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
                            {/* TODO: ウォッチリストに追加するアクションを実装すること */}
                            <div className="col-md-6 text-center">
                              <button
                                className="btn btn-outline-warning"
                                onClick={() => addWatchList(movie)}
                              >
                                Watch List
                              </button>
                            </div>
                            {/* TODO: レビューを登録するアクションを実装すること。ただし、既にレビュー済みであればこのリンクは消すこと */}
                            <div className="col-md-6 text-center">
                              <Link
                                className="btn btn-outline-success"
                                href={`/movie/${movie.id}/review/create`}
                              >
                                Review
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
      )}
    </div>
  );
}
