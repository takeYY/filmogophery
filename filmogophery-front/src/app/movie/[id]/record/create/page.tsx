"use client";

import React, { useEffect, useState } from "react";
import { WatchMedia, MovieDetail, Movie } from "@/interface/movie";
import Image from "next/image";

export default function Page({ params }: { params: { id: string } }) {
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const [watchMedia, setMedia] = useState<WatchMedia[]>([]);
  const [movie, setMovie] = useState<MovieDetail>();

  useEffect(() => {
    const media = fetchMedia();
    const mov = fetchMovie(params.id);

    setMedia(media);
    setMovie(mov);
  }, [params.id]);

  //const watchMedia = fetchMedia();
  //const movie = fetchMovie(params.id);

  const [rangeValue, onChange] = useState(movie?.impression.rating.toString());

  // /*
  useEffect(() => {
    const tooltip = document.getElementById("rangeValue");
    console.log(tooltip);
  });
  // */

  async function onSubmit(formData: FormData) {
    setIsLoading(true);
    setError(null);

    try {
      /*  API との通信
      const response = await fetch("/api/submit", {
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        throw new Error("Failed to submit the data. Please try again.");
      }

      const data = await response.json();
       */

      console.log(`media: ` + formData.get("media"));
      console.log(`watchDate: ` + formData.get("watch_date"));
      console.log(`rating: ` + formData.get("rating"));
      console.log(`note: ` + formData.get("note"));
    } catch (error) {
      // Capture the error message to display to the user
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  }

  if (!movie) {
    return <div></div>;
  }

  return (
    <div className="container-fluid pb-4">
      <h3 className="text-center mb-4">Create Movie Watch Record</h3>

      <form action={onSubmit}>
        <div className="my-4 row">
          <div className="col-md-4">
            <Image
              src={movie.posterURL}
              className="img-fluid"
              width={300}
              height={300}
              alt="ポスター"
            ></Image>
          </div>

          <div className="col-md-8">
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
                  {watchMedia.map((wm: WatchMedia, i: number) => {
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
                        <label className="form-check-label" htmlFor={wm.code}>
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
                  name="watch_date"
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
                  rangeValue === undefined
                    ? movie?.impression.rating
                    : rangeValue
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
                    movie.impression.note === undefined
                      ? ""
                      : movie.impression.note
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
      </form>
    </div>
  );
}

function fetchMedia(): WatchMedia[] {
  console.log("データ取得中...");
  const watchMedia: WatchMedia[] = [
    { code: "amazon_prime", name: "Prime Video" },
    { code: "netflix", name: "Netflix" },
    { code: "u_next", name: "U-NEXT" },
    { code: "disney_plus", name: "Disney+" },
    { code: "youtube", name: "YouTube" },
    { code: "apple_tv", name: "Apple TV+" },
    { code: "hulu", name: "Hulu" },
    { code: "d_anime", name: "dアニメ" },
    { code: "telasa", name: "TELASA" },
    { code: "cinema", name: "映画館" },
    { code: "unknown", name: "不明" },
  ];
  console.log("データ取得: 完了");

  return watchMedia;
}

function fetchMovie(id: string): MovieDetail {
  try {
    console.log("データ取得中...");
    const movieDetail: MovieDetail = {
      id: parseInt(id),
      title: "ターミネーター",
      overview:
        "アメリカのとある街、深夜突如奇怪な放電と共に屈強な肉体をもった男が現れる。同じく...",
      release_date: new Date("1985-05-04"),
      run_time: 108,
      genres: ["アクション", "SF"],
      posterURL:
        "https://image.tmdb.org/t/p/w600_and_h900_bestv2/iK2YBfD7DdaNXZALQhhT9uTN9Rc.jpg",
      vote_average: 3.85,
      vote_count: 12832,
      series: {
        name: "ターミネーターシリーズ",
        posterURL:
          "https://image.tmdb.org/t/p/w600_and_h900_bestv2/pF5GIijY2fyZcByqNDzhS8v4h1x.jpg",
      },
      impression: {
        status: "鑑賞済み",
        rating: 4.3,
        note: "ターミネーターの元祖という感じで、恐ろしさと希望が織り成す圧巻の作品。今観るとCGのぎこちなさが目立つが、それが逆に怖さを演出している。",
      },
      watchRecords: [
        {
          watch_date: new Date("2024-08-11"),
          watch_media: "U-NEXT",
        },
        {
          watch_date: new Date("2024-02-03"),
          watch_media: "Netflix",
        },
        {
          watch_date: new Date("2024-01-02"),
          watch_media: "Amazon Prime Video",
        },
      ],
    };
    console.log("データ取得: 完了!");
    return movieDetail;
  } catch (error) {
    console.log("データ取得エラー");
    return {
      id: parseInt(id),
      title: "ターミネーター",
      overview:
        "アメリカのとある街、深夜突如奇怪な放電と共に屈強な肉体をもった男が現れる。同じく...",
      release_date: new Date("1985-05-04"),
      run_time: 108,
      genres: ["アクション", "SF"],
      posterURL:
        "https://image.tmdb.org/t/p/w600_and_h900_bestv2/iK2YBfD7DdaNXZALQhhT9uTN9Rc.jpg",
      vote_average: 3.85,
      vote_count: 12832,
      series: {
        name: "ターミネーターシリーズ",
        posterURL:
          "https://image.tmdb.org/t/p/w600_and_h900_bestv2/pF5GIijY2fyZcByqNDzhS8v4h1x.jpg",
      },
      impression: {
        status: "鑑賞済み",
        rating: 4.3,
        note: "初代です！",
      },
      watchRecords: [
        {
          watch_date: new Date("2024-08-11"),
          watch_media: "U-NEXT",
        },
        {
          watch_date: new Date("2024-02-03"),
          watch_media: "Netflix",
        },
        {
          watch_date: new Date("2024-01-02"),
          watch_media: "Amazon Prime Video",
        },
      ],
    };
  }
}
