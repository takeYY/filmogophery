// watchlist/page.tsx
/**
 * ウォッチリスト
 * パス: /watchlist
 */
"use client";

import { posterUrlPrefix } from "@/constants/poster";
import { useAuth } from "@/hooks/useAuth";
import { Genre, Watchlist } from "@/interface/index";
import Image from "next/image";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function Page() {
  const searchParams = useSearchParams();
  const query = searchParams.get("query");

  const [loading, setLoading] = useState(true);
  const [watchlist, setWatchlist] = useState<Watchlist[]>([]);

  const token = useAuth();
  const accessToken = token ? token.accessToken : null;

  const headers: HeadersInit = {};
  if (accessToken) {
    headers.Authorization = `Bearer ${accessToken}`;
  }

  useEffect(() => {
    const fetchWatchlist = async () => {
      setLoading(true);
      try {
        const response = await fetch(`/api/watchlist`, {
          method: "GET",
          headers,
        });
        const watchlist: Watchlist[] = await response.json();
        console.log("watchlistのデータ取得: 完了");
        console.log("%o", watchlist);

        setWatchlist(watchlist);
      } catch {
        console.log("watchlistのデータ取得: エラー。空配列で定義します");
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
          {watchlist.map((wl: Watchlist) => {
            return (
              <div className="col" key={wl.id}>
                <div className="card mb-2 bg-dark border-info position-relative">
                  <div className="row g-0">
                    <div className="col-md-4">
                      {/* ポスター */}
                      <Image
                        src={
                          posterUrlPrefix +
                          (wl.movie.posterURL
                            ? wl.movie.posterURL
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
                        <h5 className="card-title">{wl.movie.title}</h5>
                        {/* ジャンル */}
                        {wl.movie.genres.length !== 0 && (
                          <div className="card-text d-grid gap-2 d-md-block">
                            {wl.movie.genres.map((g: Genre, i: number) => {
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
                        <p className="card-text">
                          公開日：{wl.movie.releaseDate}
                        </p>
                        {/* 概要 */}
                        <p className="card-text">
                          {wl.movie.overview.length > 40
                            ? wl.movie.overview.substring(0, 37) + "..."
                            : wl.movie.overview}
                        </p>

                        <div className="border-top border-success">
                          <div className="row mt-2">
                            {/* TODO: レビューを登録するアクションを実装すること。ただし、既にレビュー済みであればこのリンクは消すこと */}
                            <div className="col text-center">
                              <Link
                                className="btn btn-outline-success"
                                href={`/movie/${wl.movie.id}/review/create`}
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
