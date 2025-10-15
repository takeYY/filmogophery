"use client";

import React, { useEffect, useState } from "react";
import { WatchMedia, MovieDetail } from "@/interface/movie";
import StarRating from "@/app/components/Rating";
import Image from "next/image";
import { posterUrlPrefix } from "@/constants/poster";
import { useRouter } from "next/navigation";

// 感傷履歴を作るページ
export default function Page({ params }: { params: { id: string } }) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [watchMedia, setMedia] = useState<WatchMedia[]>();
  const [movieDetail, setMovie] = useState<MovieDetail>();
  // TODO: rangeValue に値が入ると、"" になってしまうので直したい...
  const [rangeValue, onChange] = useState<string>(
    movieDetail?.impression?.rating?.toString()
      ? movieDetail?.impression.rating?.toString()
      : ""
  );

  useEffect(() => {
    const fetchMedia = async () => {
      console.log("mediaのデータ取得中...");
      try {
        const response = await fetch(`/api/media`, { method: "GET" });
        const media: WatchMedia[] = await response.json();

        console.log("mediaのデータ取得: 完了");

        return setMedia(media);
      } catch {
        console.log("mediaデータ取得エラー");
      }
    };

    const fetchMovie = async () => {
      console.log("movieDetailのデータ取得中...");
      try {
        const response = await fetch(`/api/movie?id=${params.id}`, {
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

    fetchMedia();
    fetchMovie();
  }, [params.id]);

  async function onSubmit(formData: FormData) {
    setIsLoading(true);
    try {
      const jsonData = {
        mediaCode: formData.get("mediaCode"),
        date: formData.get("date"),
      };
      console.log("page payload:", jsonData);
      const response = await fetch(`/api/movies/${params.id}/records`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(jsonData),
      });

      if (!response.ok) {
        throw new Error("Failed to submit the data. Please try again.");
      }

      router.push(`/movie/${params.id}?updated=true`);
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
            movieDetail.impression?.status === "鑑賞済み"
              ? "border-success"
              : ""
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
                {/* 鑑賞媒体 */}
                <div className="form-group row">
                  <label className="col-sm-4 col-form-label d-flex align-items-center">
                    <div className="bg-transparent badge border border-danger rounded-pill">
                      必須
                    </div>
                    {"　"}
                    鑑賞媒体
                  </label>
                  <div className="col-sm-8 d-flex align-items-center">
                    <div className="px-2 row">
                      {watchMedia !== undefined &&
                        watchMedia.map((wm: WatchMedia, i: number) => {
                          return (
                            <div
                              key={i}
                              className="form-check form-check-inline col-md-3"
                            >
                              <input
                                className="form-check-input"
                                type="radio"
                                name="media"
                                id={wm.code}
                                value={wm.code}
                              />
                              <label
                                className="form-check-label"
                                htmlFor={wm.code}
                              >
                                {wm.name}
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
                      name="watchDate"
                      defaultValue={new Date().toLocaleDateString("sv-SE")}
                    />
                  </div>
                </div>

                <div className="h4 pb-2 mb-4 text-success border-bottom border-success mt-4">
                  My Impression
                </div>

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
                        movieDetail?.impression?.note
                          ? movieDetail.impression.note
                          : ""
                      }`}
                    />
                  </div>
                </div>

                <div className="text-center mt-4">
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
