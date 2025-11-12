"use client";

import React, { useEffect, useState } from "react";
import { Platform, MovieDetail, Genre } from "@/interface/index";
import StarRating from "@/app/components/Rating";
import Image from "next/image";
import { posterUrlPrefix } from "@/constants/poster";
import { useRouter } from "next/navigation";

export default function Page({
  params,
}: {
  params: { id: string; reviewId: string };
}) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [platforms, setPlatforms] = useState<Platform[]>();
  const [movieDetail, setMovie] = useState<MovieDetail>();
  const [rangeValue, onChange] = useState<string>("");

  // movieDetailが更新されたときにrangeValueを設定
  useEffect(() => {
    if (movieDetail?.review?.rating) {
      onChange(movieDetail.review.rating.toString());
    }
  }, [movieDetail]);

  useEffect(() => {
    const fetchPlatforms = async () => {
      console.log("platformsのデータ取得中...");
      try {
        const response = await fetch(`/api/platforms`, { method: "GET" });
        const platforms: Platform[] = await response.json();

        console.log("platformsのデータ取得: 完了");

        return setPlatforms(platforms);
      } catch {
        console.log("platformsデータ取得エラー");
      }
    };

    const fetchMovie = async () => {
      console.log("movieDetailのデータ取得中...");
      try {
        const response = await fetch(`/api/movies/${params.id}`, {
          method: "GET",
        });
        const movieDetail: MovieDetail = await response.json();

        console.log("movieDetailのデータ取得: 完了");
        console.log("%o", movieDetail);

        return setMovie(movieDetail);
      } catch {
        console.log("moviesのデータ取得: エラー。空配列で定義します");
        return setMovie(undefined);
      }
    };

    fetchPlatforms();
    fetchMovie();
  }, [params.id]);

  async function onSubmit(formData: FormData) {
    setIsLoading(true);
    try {
      const jsonData = {
        platformId: formData.get("platformId"),
        watchedDate: formData.get("watchedDate"),
      };
      console.log("page payload:", jsonData);
      const response = await fetch(
        `/api/movies/${params.id}/reviews/${movieDetail?.review?.id}`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(jsonData),
        }
      );

      if (!response.ok) {
        throw new Error("Failed to submit the data. Please try again.");
      }

      router.push(`/movie/${params.id}?updated=true&t=${Date.now()}`);
      router.refresh();
    } catch (error) {
      // Capture the error message to display to the user
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  }

  if (!movieDetail) {
    return <div></div>;
  }

  return (
    <div className="container-fluid pb-4">
      <h3 className="text-center mb-4">Create Movie Watch Record</h3>

      <form action={onSubmit}>
        <div
          className={`card mb-3 bg-dark ${
            movieDetail.review !== null ? "border-success" : ""
          }`}
        >
          <div className="row g-0">
            <div className="col-md-3">
              {/* ポスター */}
              <Image
                src={
                  posterUrlPrefix +
                  (movieDetail.posterURL
                    ? movieDetail.posterURL
                    : "/Agz71U0wcesx87micVn731Z1vPu.jpg")
                }
                className="img-fluid rounded-start"
                width={350}
                height={350}
                alt="ポスター"
                priority={false}
              />

              {/* 一般の評価 */}
              <div className="justify-content-center">
                <StarRating
                  rating={movieDetail.voteAverage}
                  size={20}
                  starColor={"#0dcaf0"}
                  sumReview={movieDetail.voteCount.toString()}
                />
              </div>
            </div>

            <div className="col-md-9">
              <div className="card-body text-light">
                {/* タイトル */}
                <h5 className="card-title">{movieDetail.title}</h5>
                {/* ジャンル */}
                {movieDetail.genres.length !== 0 && (
                  <div className="card-text d-grid gap-2 d-md-block">
                    {movieDetail.genres.map((g: Genre, i: number) => {
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
                <p className="card-text">公開日：{movieDetail.releaseDate}</p>
                {/* 上映時間 */}
                <p className="card-text">
                  上映時間：{movieDetail.runtimeMinutes}分
                </p>
                {/* プラットフォーム */}
                <div className="form-group row">
                  <label className="col-sm-4 col-form-label d-flex align-items-center">
                    <div className="bg-transparent badge border border-danger rounded-pill">
                      必須
                    </div>
                    {"　"}
                    プラットフォーム
                  </label>
                  <div className="col-sm-8 d-flex align-items-center">
                    <div className="px-2 row">
                      {platforms !== undefined &&
                        platforms.map((p: Platform, i: number) => {
                          return (
                            <div
                              key={i}
                              className="form-check form-check-inline col-md-3"
                            >
                              <input
                                className="form-check-input"
                                type="radio"
                                name="platformId"
                                id={p.code}
                                value={i + 1}
                              />
                              <label
                                className="form-check-label"
                                htmlFor={p.code}
                              >
                                {p.name}
                              </label>
                            </div>
                          );
                        })}
                    </div>
                  </div>
                </div>

                {/* 鑑賞日 */}
                <div className="form-group row mt-4">
                  <label className="col-sm-4 col-form-label d-flex align-items-center">
                    <div className="bg-transparent badge border border-info rounded-pill">
                      任意
                    </div>
                    {"　"}鑑賞日
                  </label>
                  <div className="col-sm-8">
                    <input
                      type="date"
                      className="form-control w-50 bg-dark text-light"
                      name="watchedDate"
                      defaultValue={new Date().toLocaleDateString("sv-SE")}
                    />
                  </div>
                </div>

                <div className="h4 pb-2 mb-4 text-success border-bottom border-success mt-4">
                  Review
                </div>

                {/* TODO: レビュー内容が前回と違えば更新するようにAPIリクエストすること!! */}
                <div className="form-group row mt-4">
                  <label className="col-sm-4 col-form-label d-flex align-items-center">
                    <div className="bg-transparent badge border border-info rounded-pill">
                      任意
                    </div>
                    {"　"}評価
                  </label>
                  <div className="col-sm-8">
                    <div id="rangeValue">{`${
                      rangeValue ? rangeValue : "評価なし"
                    }`}</div>
                    <input
                      type="range"
                      name="rating"
                      step={0.1}
                      max={5.0}
                      min={1.0}
                      value={rangeValue}
                      className="form-range"
                      onChange={({ target: { value: radius } }) => {
                        onChange(radius);
                      }}
                    />
                  </div>
                </div>

                {/* 感想 */}
                <div className="form-group row mt-4">
                  <label className="col-sm-4 col-form-label d-flex align-items-center">
                    <div className="bg-transparent badge border border-info rounded-pill">
                      任意
                    </div>
                    {"　"}感想
                  </label>
                  <div className="col-sm-8">
                    <textarea
                      className="form-control bg-dark text-light"
                      name="note"
                      defaultValue={`${
                        movieDetail?.review?.comment
                          ? movieDetail.review.comment
                          : ""
                      }`}
                    />
                  </div>
                </div>

                <div className="text-center mt-4">
                  <button
                    type="button"
                    onClick={() => router.back()}
                    className="btn btn-outline-light me-3"
                  >
                    Cancel
                  </button>

                  <button
                    type="submit"
                    disabled={isLoading}
                    className="btn btn-outline-primary"
                  >
                    {isLoading ? "Loading..." : "Create"}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
}
