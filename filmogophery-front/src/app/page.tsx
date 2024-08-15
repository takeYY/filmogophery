"use client";

import React, { useEffect, useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { Movie, Genre } from "@/interface/movie";
import { posterUrlPrefix } from "@/constants/poster";

export default function Home() {
  const router = useRouter();
  const [movies, setMovies] = useState<Movie[]>();

  useEffect(() => {
    const fetchMovies = async () => {
      try {
        const response = await fetch(`/api/movies`, { method: "GET" });
        const movies: Movie[] = await response.json();
        console.log("moviesのデータ取得: 完了");

        return setMovies(movies);
      } catch {
        console.log("moviesのデータ取得: エラー。空配列で定義します");
        return setMovies([]);
      }
    };

    fetchMovies();
  }, []);

  if (!movies) {
    return <div></div>;
  }
  return (
    <main>
      <div className="container pb-4">
        <h3 className="text-center mb-4">Home</h3>

        {/* TODO: レイアウトが崩れているので、直すこと!! */}
        <div className="card-columns">
          {movies &&
            movies.map((movie: Movie, index: number) => {
              return (
                <button
                  key={index}
                  className="card bg-dark"
                  onClick={() => router.push(`/movie/${movie.id}`)}
                >
                  <div className="row no-gutters">
                    {/* ポスター */}
                    <div className="col-md-4">
                      <Image
                        src={
                          posterUrlPrefix +
                          (movie.poster
                            ? movie.poster.url
                            : "/Agz71U0wcesx87micVn731Z1vPu.jpg")
                        }
                        alt="ポスター画像"
                        className="img-fluid"
                        width={75}
                        height={75}
                      />
                    </div>
                    <div className="col-md-8">
                      <div className="card-body">
                        {/* タイトル */}
                        <h5 className="card-title text-light">{movie.title}</h5>
                        {/* 概要 */}
                        <p className="card-text text-light">{movie.overview}</p>
                        {/* ジャンル */}
                        {movie.genres &&
                          movie.genres.map((genre: Genre, i: number) => {
                            return (
                              <p key={i} className="badge text-bg-secondary">
                                {genre.name}
                              </p>
                            );
                          })}
                      </div>
                    </div>
                  </div>
                </button>
              );
            })}
        </div>
      </div>
    </main>
  );
}
