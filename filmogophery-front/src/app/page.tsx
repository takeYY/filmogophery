"use client";

import React, { useEffect, useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { Movie, Genre } from "@/interface/movie";
import { posterUrlPrefix } from "@/constants/poster";
import { Carousel } from "react-bootstrap";

export default function Home() {
  const router = useRouter();
  const [movies, setMovies] = useState<Movie[]>();

  function separatedMovie(movies: Movie[] | undefined, size: number) {
    if (movies?.length === 0 || movies === undefined) {
      return [[]];
    }
    return movies.flatMap((_, i, a) =>
      i % size ? [] : [movies.slice(i, i + size)]
    );
  }

  const [separated, setSeparated] = useState<Movie[][]>(
    separatedMovie(movies, 5)
  );

  useEffect(() => {
    const fetchMovies = async () => {
      try {
        const response = await fetch(`/api/movies`, { method: "GET" });
        const movies: Movie[] = await response.json();
        console.log("moviesのデータ取得: 完了");

        setMovies(movies);
        setSeparated(separatedMovie(movies, 5));
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

        <h5>最近の映画</h5>
        <Carousel pause={"hover"} className="mb-4">
          {separated.map((movies: Movie[], index: number) => {
            return (
              <Carousel.Item>
                <div className="row justify-content-md-center">
                  {movies.map((movie: Movie, i: number) => {
                    return (
                      <div className="col-md-2">
                        <Image
                          src={
                            posterUrlPrefix +
                            (movie.poster_url
                              ? movie.poster_url
                              : "/Agz71U0wcesx87micVn731Z1vPu.jpg")
                          }
                          alt="ポスター画像"
                          className="img-fluid"
                          width={200}
                          height={200}
                        />
                      </div>
                    );
                  })}
                </div>
              </Carousel.Item>
            );
          })}
        </Carousel>

        <div className="row row-cols-md-3 g-4">
          {movies.map((movie: Movie, index: number) => {
            return (
              <div className="col">
                <button
                  className="card mb-2 bg-dark border-info"
                  onClick={() => router.push(`/movie/${movie.id}`)}
                >
                  <div className="row g-0">
                    <div className="col-md-4">
                      {/* ポスター */}
                      <Image
                        src={
                          posterUrlPrefix +
                          (movie.poster_url
                            ? movie.poster_url
                            : "/Agz71U0wcesx87micVn731Z1vPu.jpg")
                        }
                        className="card-img-top"
                        alt="..."
                        width={200}
                        height={200}
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
                                  key={i}
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
                        <p className="card-text">
                          公開日：{movie.release_date.substring(0, 10)}
                        </p>
                        {/* 概要 */}
                        <p className="card-text">
                          {movie.overview.length > 40
                            ? movie.overview.substring(0, 37) + "..."
                            : movie.overview}
                        </p>
                      </div>
                    </div>
                  </div>
                </button>
              </div>
            );
          })}
        </div>
      </div>
    </main>
  );
}
