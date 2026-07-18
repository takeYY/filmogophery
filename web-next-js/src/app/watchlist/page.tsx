// watchlist/page.tsx
/**
 * ウォッチリスト
 * パス: /watchlist
 */
"use client";

import { MovieCard } from "@/components/MovieCard";
import { Watchlist } from "@/interface/index";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function Page() {
  const searchParams = useSearchParams();
  const query = searchParams.get("query");

  const [loading, setLoading] = useState(true);
  const [watchlist, setWatchlist] = useState<Watchlist[]>([]);

  useEffect(() => {
    const fetchWatchlist = async () => {
      setLoading(true);
      try {
        const response = await fetch(`/api/watchlist`, {
          method: "GET",
        });
        const data: Watchlist[] = await response.json();
        setWatchlist(data);
      } catch {
        setWatchlist([]);
      } finally {
        setLoading(false);
      }
    };
    fetchWatchlist();
  }, [query]);

  return (
    <div className="container pb-4">
      <h3 className="text-center mb-4">Watchlist</h3>

      {loading ? (
        <div className="text-center">
          <div className="spinner-border text-info" role="status">
            <span className="visually-hidden">Loading...</span>
          </div>
        </div>
      ) : watchlist.length === 0 ? (
        <div className="text-center text-light">
          <p>ウォッチリストはありませんでした。</p>
        </div>
      ) : (
        <div className="row row-cols-md-3 g-4">
          {watchlist.map((wl: Watchlist) => (
            <div className="col" key={wl.id}>
              <MovieCard
                movie={wl.movie}
                imageWidth={250}
                imageHeight={375}
                actions={
                  // TODO: レビュー済みであればこのリンクは消すこと
                  <Link
                    className="btn btn-outline-success"
                    href={`/movie/${wl.movie.id}/review/create`}
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
