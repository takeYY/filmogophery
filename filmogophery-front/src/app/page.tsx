"use client";

import React, { useEffect, useState } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { MovieNeo, Genre } from "@/interface/movie";
import { posterUrlPrefix } from "@/constants/poster";
import { Carousel } from "react-bootstrap";

export default function Home() {
  const router = useRouter();
  const [movies, setMovies] = useState<MovieNeo[]>();
  const [isLoading, setIsLoading] = useState(true);

  function separatedMovie(movies: MovieNeo[] | undefined, size: number) {
    if (!movies || movies.length === 0) {
      return [[]];
    }
    const result = [];
    for (let i = 0; i < movies.length; i += size) {
      result.push(movies.slice(i, i + size));
    }
    return result;
  }

  const [separated, setSeparated] = useState<MovieNeo[][]>([[]]);

  useEffect(() => {
    const fetchMovies = async () => {
      setIsLoading(true);
      try {
        const response = await fetch(`/api/movies`, { method: "GET" });
        const movies: MovieNeo[] = await response.json();
        console.log("moviesのデータ取得: 完了");

        setMovies(movies);
        setSeparated(separatedMovie(movies, 5));
      } catch (error) {
        console.log("moviesのデータ取得: エラー。空配列で定義します");
        setMovies([]);
      } finally {
        setIsLoading(false);
      }
    };

    fetchMovies();
  }, []);

  if (isLoading) {
    return (
      <div
        className="container d-flex justify-content-center align-items-center"
        style={{ minHeight: "50vh" }}
      >
        <div className="spinner-border text-info" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    );
  }

  if (!movies || movies.length === 0) {
    return <div className="container text-center py-5">No movies found</div>;
  }

  return (
    <main>
      <div className="container pb-4">
        <h3 className="text-center mb-4">Home</h3>

        <h5>最近の映画</h5>
        <Carousel pause={"hover"} className="mb-4">
          {separated.map((movies: MovieNeo[], index: number) => {
            return (
              <Carousel.Item key={`carousel-item-${index}`}>
                <div className="row justify-content-md-center">
                  {movies.map((movie: MovieNeo, i: number) => {
                    return (
                      <div
                        className="col-md-2"
                        key={`carousel-movie-${movie.id || i}`}
                      >
                        <Image
                          src={
                            posterUrlPrefix +
                            (movie.posterURL
                              ? movie.posterURL
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
          {movies.map((movie: MovieNeo, index: number) => {
            return (
              <div className="col" key={`movie-card-${movie.id || index}`}>
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
                          (movie.posterURL
                            ? movie.posterURL
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
                          公開日：{movie.releaseDate.substring(0, 10)}
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
