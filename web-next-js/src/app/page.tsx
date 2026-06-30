// app/page.tsx
/**
 * ホームページ
 * パス: /
 */

"use client";

import { posterUrlPrefix } from "@/constants/poster";
import { useAuth } from "@/hooks/useAuth";
import { Movie, TrendingMovie } from "@/interface/index";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { Carousel } from "react-bootstrap";

export default function Home() {
  const router = useRouter();
  const { checked } = useAuth();

  const [movies, setMovies] = useState<Movie[]>();
  const [trending, setTrending] = useState<TrendingMovie[]>();
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const observerTarget = useRef<HTMLDivElement>(null);

  function separateTrending(
    trending: TrendingMovie[] | undefined,
    size: number,
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
    if (!checked) return;
    const loadInitialData = async () => {
      setIsLoading(true);
      try {
        const initialMovies = await fetchMovies(0);
        const res = await fetch(`/api/trending/movies`, {
          method: "GET",
        });
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
  }, [checked]);

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
      { threshold: 0.1 },
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
                  {trending.map((trend: TrendingMovie, i: number) => {
                    return (
                      <div
                        className="col-md-2"
                        key={`carousel-movie-${trend.id || i}`}
                      >
                        <Image
                          src={
                            posterUrlPrefix +
                            (trend.posterURL
                              ? trend.posterURL
                              : "/Agz71U0wcesx87micVn731Z1vPu.jpg")
                          }
                          alt="ポスター画像"
                          className="img-fluid"
                          width={200}
                          height={200}
                          onClick={() =>
                            router.push(
                              trend.hasReview
                                ? `/movie/${trend.id}`
                                : `/movie/${trend.id}/review/create`,
                            )
                          }
                          style={{ cursor: "pointer" }}
                        />
                      </div>
                    );
                  })}
                </div>
              </Carousel.Item>
            );
          })}
        </Carousel>

        <h5>レビュー済み映画</h5>
        {!movies || movies.length === 0 ? (
          <div className="text-center py-4">
            <p className="text-muted mb-3">
              まだレビューした映画がありません。
            </p>
            <button
              className="btn btn-outline-info"
              onClick={() => router.push("/search")}
            >
              映画を探す
            </button>
          </div>
        ) : (
          <>
            <div className="row row-cols-md-3 g-4">
              {movies.map((movie: Movie, index: number) => (
                <div className="col" key={`movie-card-${movie.id || index}`}>
                  <MovieCard
                    movie={movie}
                    onClick={() => router.push(`/movie/${movie.id}`)}
                  />
                </div>
              ))}
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
          </>
        )}
      </div>
    </main>
  );
}
