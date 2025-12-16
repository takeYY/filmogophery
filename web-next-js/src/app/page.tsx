// app/page.tsx
/**
 * ホームページ
 * パス: /
 */

"use client";

import React, { useEffect, useState, useRef } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { Movie, Genre, TrendingMovie } from "@/interface/index";
import { posterUrlPrefix } from "@/constants/poster";
import { Carousel } from "react-bootstrap";

export default function Home() {
  const router = useRouter();
  const [movies, setMovies] = useState<Movie[]>();
  const [trending, setTrending] = useState<TrendingMovie[]>();
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const observerTarget = useRef<HTMLDivElement>(null);

  function separateTrending(
    trending: TrendingMovie[] | undefined,
    size: number
  ) {
    if (!trending || trending.length === 0) {
      return [[]];
    }
    const result = [];
    for (let i = 0; i < trending.length; i += size) {
      result.push(trending.slice(i, i + size));
    }
    return result;
  }

  const [separated, setSeparated] = useState<TrendingMovie[][]>([[]]);

  const fetchMovies = async (currentOffset: number) => {
    try {
      const response = await fetch(`/api/movies?offset=${currentOffset}`, {
        method: "GET",
      });
      const newMovies: Movie[] = await response.json();
      return newMovies;
    } catch (error) {
      console.log("moviesのデータ取得: エラー");
      return [];
    }
  };

  useEffect(() => {
    const loadInitialData = async () => {
      setIsLoading(true);
      try {
        const initialMovies = await fetchMovies(0);
        const res = await fetch(`/api/trending/movies`, { method: "GET" });
        const trending: TrendingMovie[] = await res.json();
        console.log("moviesのデータ取得: 完了");

        setMovies(initialMovies);
        setTrending(trending);
        setSeparated(separateTrending(trending, 5));
        setOffset(12);
        setHasMore(initialMovies.length === 12);
      } catch (error) {
        console.log("moviesのデータ取得: エラー。空配列で定義します");
        setMovies([]);
      } finally {
        setIsLoading(false);
      }
    };

    loadInitialData();
  }, []);

  useEffect(() => {
    const observer = new IntersectionObserver(
      async (entries) => {
        if (
          entries[0].isIntersecting &&
          !isLoadingMore &&
          hasMore &&
          movies &&
          movies.length > 0
        ) {
          setIsLoadingMore(true);
          const newMovies = await fetchMovies(offset);
          if (newMovies.length > 0) {
            setMovies((prev) => [...(prev || []), ...newMovies]);
            setOffset((prev) => prev + 12);
          } else {
            setHasMore(false);
          }
          setIsLoadingMore(false);
        }
      },
      { threshold: 0.1 }
    );

    if (observerTarget.current) {
      observer.observe(observerTarget.current);
    }

    return () => observer.disconnect();
  }, [offset, isLoadingMore, hasMore, movies?.length]);

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
          {separated.map((trending: TrendingMovie[], index: number) => {
            return (
              <Carousel.Item key={`carousel-item-${index}`}>
                <div className="row justify-content-md-center">
                  {trending.map((t: TrendingMovie, i: number) => {
                    return (
                      <div
                        className="col-md-2"
                        key={`carousel-movie-${t.id || i}`}
                      >
                        <Image
                          src={
                            posterUrlPrefix +
                            (t.posterURL
                              ? t.posterURL
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

        <div
          ref={observerTarget}
          style={{ height: "20px", marginTop: "20px" }}
        />
        {isLoadingMore && (
          <div className="text-center py-3">
            <div className="spinner-border text-info" role="status">
              <span className="visually-hidden">Loading...</span>
            </div>
          </div>
        )}
      </div>
    </main>
  );
}
