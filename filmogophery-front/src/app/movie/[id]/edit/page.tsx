"use client";

import React, { useEffect, useState } from "react";
import { Impression } from "@/interface/movie";
import { impressionData } from "@/app/lib/movies_detail_data";
import Image from "next/image";
import StarRating from "@/app/components/Rating";
import { MovieDetail } from "@/interface/movie";
import { movieDetailData } from "@/app/lib/movies_detail_data";
import { posterUrlPrefix } from "@/constants/poster";

// 感想を編集するページ
export default function Page({ params }: { params: { id: string } }) {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [impression, setImpression] = useState<Impression>();
  const [movieDetail, setMovie] = useState<MovieDetail>();
  const [rangeValue, onChange] = useState<string>(
    movieDetail?.impression?.rating?.toString()
      ? movieDetail?.impression.rating?.toString()
      : ""
  );

  useEffect(() => {
    const fetchImpression = async () => {
      console.log("impressionのデータ取得中...");
      try {
        // FIXME: APIから取得するように修正すること
        const impression: Impression | undefined = impressionData.get(
          parseInt(params.id, 10)
        );
        console.log("impressionのデータ: 完了");
        console.log("%o", impression);

        return setImpression(impression);
      } catch {
        console.log("impression取得エラー");
        return setImpression(undefined);
      }
    };
    const fetchMovieDetail = async () => {
      console.log("movieDetailのデータ取得中...");
      try {
        // const response = await fetch(`/api/movie?id=${params.id}`, {
        //   method: "GET",
        // });
        // const movie: MovieDetail = await response.json();

        // FIXME: APIから取得するように修正すること
        const movieDetail: MovieDetail | undefined = movieDetailData.get(
          params.id
        );
        console.log("movieDetailのデータ取得: 完了");
        console.log("%o", movieDetail);

        return setMovie(movieDetail);
      } catch {
        console.log("moviesのデータ取得: エラー。空配列で定義します");
        return setMovie(undefined);
      }
    };

    fetchImpression();
    fetchMovieDetail();
  }, [params.id]);

  async function onSubmit(formData: FormData) {
    setIsLoading(true);
    try {
      // TODO: API との通信
    } catch (error) {
      console.log(error);
    } finally {
      setIsLoading(false);
    }
  }

  if (!movieDetail) {
    return <div></div>;
  }

  return (
    <div className="container-fluid pb-4">
      <h3 className="text-center mb-4">Edit Movie Impression</h3>

      <form action={onSubmit}>
        <div
          className={`card mb-3 bg-dark ${
            impression?.status === "鑑賞済み" ? "border-success" : ""
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
                {/* 評価 */}
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
                    className="btn btn-outline-success"
                  >
                    {isLoading ? "Loading..." : "Edit"}
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
