"use client";

import React, { useEffect, useState } from "react";
import StarRating from "@/app/components/Rating";
import { MovieDetail, WatchRecord } from "@/interface/movie";
import Image from "next/image";
import Link from "next/link";
import { posterUrlPrefix } from "@/constants/poster";
import { movieDetailData } from "@/app/lib/movies_detail_data";

export default function Page({ params }: { params: { id: string } }) {
  const [movie, setMovie] = useState<MovieDetail | null>();

  useEffect(() => {
    const fetchMovies = async () => {
      console.log("movieのデータ取得中...");
      try {
        // const response = await fetch(`/api/movie?id=${params.id}`, {
        //   method: "GET",
        // });
        // const movie: MovieDetail = await response.json();
        const movie: MovieDetail | null =
          movieDetailData.get(params.id) ?? null;
        console.log("moviesのデータ取得: 完了");
        console.log("%o", movie);

        return setMovie(movie);
      } catch {
        console.log("moviesのデータ取得: エラー。空配列で定義します");
        return setMovie(undefined);
      }
    };

    fetchMovies();
  }, [params.id]);

  if (!movie) {
    return <div>Movie({params.id}) is not found</div>;
  }

  // const movie = fetchMovie(params.id);
  return (
    <div className="container-fluid pb-4">
      <h3 className="text-center mb-4">Movie Detail</h3>

      <div
        className={`card mb-3 bg-dark ${
          movie.impression?.status === "鑑賞済み" ? "border-success" : ""
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
              {movie.impression?.rating && (
                <div className="card-text">
                  <StarRating rating={movie.impression.rating} size={20} />
                </div>
              )}
              {/* ジャンル */}
              {movie.genres.length !== 0 && (
                <div className="card-text d-grid gap-2 d-md-block">
                  {movie.genres.map((g: string, i: number) => {
                    return (
                      <button
                        key={i}
                        type="button"
                        className="btn btn-outline-info btn-sm"
                      >
                        {g}
                      </button>
                    );
                  })}
                </div>
              )}
              {/* 公開日 */}
              <p className="card-text">公開日：{movie.releaseDate}</p>
              {/* 上映時間 */}
              <p className="card-text">上映時間：{movie.runTime}分</p>
              {/* 概要 */}
              <p className="card-text">{movie.overview}</p>
              {/* 感想 */}
              {movie.impression?.note && (
                <div className="p-3 bg-success bg-opacity-10 border border-success border-start-0 border-end-0">
                  {movie.impression?.note}
                </div>
              )}
              {/* */}
            </div>
            {/* 視聴履歴 */}
            <div className="card-footer border-success text-light">
              <div>視聴履歴</div>
              {!movie.watchRecords.length && <div>なし</div>}

              {movie.watchRecords.length !== 0 && (
                <dl className="row">
                  {movie.watchRecords.map((r: WatchRecord, i: number) => {
                    return (
                      <div key={i}>
                        <dt className="col-md-1 bg-transparent badge border border-primary rounded-pill">
                          {`${calcDiffDate(new Date(r.watchDate))}日前`}
                        </dt>
                        <dd className="col-md-10">
                          <dl className="row">
                            <dt className="col-md-4">{r.watchDate}</dt>
                            <dd className="col-md-8">{r.watchMedia}</dd>
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
            href={`/movie/${params.id}/record/create`}
          >
            視聴履歴を作成
          </Link>
        </div>
        <div className="col-md-2"></div>
      </div>
    </div>
  );
}

function calcDiffDate(target: Date): string {
  const now = new Date();
  const diff = now.getTime() - target.getTime();
  const result = Math.ceil(diff / (1000 * 60 * 60 * 24));

  return Math.abs(result).toString();
}
