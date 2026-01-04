// app/movie/[id]/page.tsx
/**
 * 映画一覧ページ
 * パス: /movie/[id]
 */

"use client";

import StarRating from "@/components/Rating";
import { posterUrlPrefix } from "@/constants/poster";
import { useAuth } from "@/hooks/useAuth";
import { Genre, MovieDetail, WatchHistory } from "@/interface/index";
import { formatRelativeTime } from "@/utils/date";
import Image from "next/image";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export default function Page({ params }: { params: { id: string } }) {
  const searchParams = useSearchParams();
  const isUpdated = searchParams.get("updated") === "true";
  const router = useRouter();
  const [showAlert, setShowAlert] = useState(isUpdated);
  const [movie, setMovie] = useState<MovieDetail | null>();
  const [watchHistory, setWatchHistory] = useState<WatchHistory[] | []>();

  const token = useAuth();
  const accessToken = token ? token.accessToken : null;

  const headers: HeadersInit = {};
  if (accessToken) {
    headers.Authorization = `Bearer ${accessToken}`;
  }

  // アラートの自動非表示とURL更新
  useEffect(() => {
    if (isUpdated) {
      setShowAlert(true);

      // 3秒後に自動で非表示
      const timer = setTimeout(() => {
        setShowAlert(false);
        // URLからupdatedパラメータを削除
        router.replace(`/movie/${params.id}`, { scroll: false });
      }, 3000);

      return () => clearTimeout(timer);
    }
  }, [isUpdated, params.id, router]);

  // 手動でアラートを閉じる
  const handleCloseAlert = () => {
    setShowAlert(false);
    router.replace(`/movie/${params.id}`, { scroll: false });
  };

  useEffect(() => {
    const fetchMoviesAndWatchHistory = async () => {
      console.log("movieのデータ取得中...");
      try {
        const response = await fetch(`/api/movies/${params.id}`, {
          method: "GET",
          headers,
          cache: "no-store",
        });
        const movie: MovieDetail = await response.json();
        console.log("movieのデータ取得: 完了");
        console.log("%o", movie);

        setMovie(movie);

        // 視聴履歴取得
        if (movie.review === null) {
          setWatchHistory([]);
          return;
        }

        console.log("視聴履歴を取得中...");
        try {
          const watchHistoryResponse = await fetch(
            `/api/watchHistory/${movie.review.id}`,
            {
              method: "GET",
              headers,
              cache: "no-store",
            }
          );
          const watchHistoryData: WatchHistory[] =
            await watchHistoryResponse.json();
          console.log("視聴履歴の取得: 完了");
          setWatchHistory(watchHistoryData);
        } catch {
          console.log("視聴履歴の取得: エラー。空配列で定義します");
          return setWatchHistory([]);
        }
      } catch {
        console.log("movieのデータ取得: エラー。空配列で定義します");
        setMovie(undefined);
        setWatchHistory([]);
      }
    };

    fetchMoviesAndWatchHistory();
  }, [params.id, searchParams]);

  if (!movie) {
    return <div>Movie({params.id}) is not found</div>;
  }

  // const movie = fetchMovie(params.id);
  return (
    <div className="container-fluid pb-4">
      {showAlert && (
        <div
          className="alert alert-success alert-dismissible fade show"
          role="alert"
        >
          感想の更新が完了しました！
          <button
            type="button"
            className="btn-close"
            onClick={handleCloseAlert}
          ></button>
        </div>
      )}
      <h3 className="text-center mb-4">Movie Detail</h3>

      <div
        className={`card mb-3 bg-dark ${
          movie.review !== null ? "border-success" : ""
        }`}
      >
        <div className="row g-0">
          <div className="col-md-3">
            {/* ポスター */}
            <Image
              src={
                posterUrlPrefix +
                (movie.posterURL
                  ? movie.posterURL
                  : "/Agz71U0wcesx87micVn731Z1vPu.jpg")
              }
              className="img-fluid rounded-start"
              alt="ポスター画像"
              width={350}
              height={350}
            />

            {/* 一般の評価 */}
            <div className="justify-content-center">
              <StarRating
                rating={movie.voteAverage}
                size={20}
                starColor={"#0dcaf0"}
                sumReview={movie.voteCount.toString()}
              />
            </div>
          </div>

          <div className="col-md-9">
            <div className="card-body text-light">
              {/* タイトル */}
              <h5 className="card-title">{movie.title}</h5>
              {/* 自身の評価 */}
              {movie.review?.rating && (
                <div className="card-text">
                  <StarRating rating={movie.review.rating} size={20} />
                </div>
              )}
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
              <p className="card-text">公開日：{movie.releaseDate}</p>
              {/* 上映時間 */}
              <p className="card-text">上映時間：{movie.runtimeMinutes}分</p>
              {/* 概要 */}
              <p className="card-text">{movie.overview}</p>
              {/* 感想 */}
              {movie.review?.comment && (
                <div className="p-3 bg-success bg-opacity-10 border border-success border-start-0 border-end-0">
                  {movie.review?.comment}
                </div>
              )}
              {/* */}
            </div>
            {/* 視聴履歴 */}
            {/* FIXME: 視聴履歴作成から戻ってきたときにデータが更新されない。。。 */}
            {/* TODO: 視聴履歴が多いと、ボタンの位置が下にずれていくのでここだけでスクロールできるようにしたい */}
            <div className="card-footer border-success text-light">
              <div>視聴履歴</div>
              {!watchHistory?.length && <div>なし</div>}

              {watchHistory?.length !== 0 && (
                <dl className="row">
                  {watchHistory?.map((wh: WatchHistory, i: number) => {
                    return (
                      <div key={i}>
                        <dt className="col-md-1 bg-transparent badge border border-primary rounded-pill">
                          {formatRelativeTime(new Date(wh.watchedAt))}
                        </dt>
                        <dd className="col-md-10">
                          <dl className="row">
                            <dt className="col-md-4">
                              {wh.watchedAt.split("T")[0]}
                            </dt>
                            <dd className="col-md-8">{wh.platform.name}</dd>
                          </dl>
                        </dd>
                      </div>
                    );
                  })}
                </dl>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* 編集 */}
      <div className="row text-center">
        <div className="col-md-2"></div>
        <div className="col-md-3">
          <Link
            className="btn btn-outline-success"
            href={`/movie/${params.id}/edit`}
          >
            感想を編集
          </Link>
        </div>
        <div className="col-md-2"></div>
        <div className="col-md-3">
          <Link
            className="btn btn-outline-primary"
            href={`/movie/${params.id}/review/${movie.review?.id}`}
          >
            視聴履歴を作成
          </Link>
        </div>
        <div className="col-md-2"></div>
      </div>
    </div>
  );
}
