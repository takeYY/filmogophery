"use client";

import { useAuth } from "@/hooks/useAuth";
import { Movie } from "@/interface/index";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

export default function SearchMovies() {
  const searchParams = useSearchParams();
  const query = searchParams.get("query");

  const [movies, setMovies] = useState<Movie[]>([]);
  const [loading, setLoading] = useState(true);
  const [addedToWatchlist, setAddedToWatchlist] = useState<number[]>([]);

  const token = useAuth();
  const accessToken = token ? token.accessToken : null;

  const headers: HeadersInit = {};
  if (accessToken) {
    headers.Authorization = `Bearer ${accessToken}`;
  }

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
          headers,
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

  const addWatchList = useCallback(
    async (movie: Movie) => {
      if (!accessToken) {
        console.log("未認証のため、ウォッチリストに追加できません");
        return;
      }

      try {
        const response = await fetch("/api/watchlist", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
          body: JSON.stringify({ movieId: movie.id }),
        });

        if (response.ok) {
          setAddedToWatchlist((prev) => [...prev, movie.id]);
        } else {
          console.log("failed to add watch list");
        }
      } catch {
        console.log("failed to add watch list");
      }
    },
    [accessToken],
  );

  return (
    <div className="container pb-4">
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
          {movies.map((movie: Movie) => (
            <div className="col" key={movie.id}>
              <MovieCard
                movie={movie}
                imageWidth={250}
                imageHeight={375}
                className="position-relative"
                overlay={
                  <button
                    className={`position-absolute top-0 start-0 m-2 btn btn-sm rounded-circle ${
                      addedToWatchlist.includes(movie.id)
                        ? "btn-success"
                        : "btn-warning"
                    }`}
                    onClick={() => addWatchList(movie)}
                    title="ウォッチリストに追加"
                    style={{ zIndex: 10, width: "36px", height: "36px" }}
                    disabled={addedToWatchlist.includes(movie.id)}
                  >
                    {addedToWatchlist.includes(movie.id) ? "✓" : "➕"}
                  </button>
                }
                actions={
                  // TODO: レビュー済みであればこのリンクは消すこと
                  <Link
                    className="btn btn-outline-success"
                    href={`/movie/${movie.id}/review/create`}
                  >
                    Review
                  </Link>
                }
              />
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
