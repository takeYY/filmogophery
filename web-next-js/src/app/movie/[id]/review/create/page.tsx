// app/movie/[id]/review/create/page.tsx
/**
 * 映画レビュー登録ページ
 * パス: /movie/[id]/review/create
 */

"use client";

import { PointToast } from "@/components/PointToast";
import StarRating from "@/components/Rating";
import { posterUrlPrefix } from "@/constants/poster";
import { useAuth } from "@/hooks/useAuth";
import { usePointToast } from "@/hooks/usePointToast";
import { Genre, MovieDetail, Platform } from "@/interface/index";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function Page({ params }: { params: { id: string } }) {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [movieDetail, setMovie] = useState<MovieDetail>();
  const [rangeValue, onChange] = useState<string>("");
  const [platforms, setPlatforms] = useState<Platform[]>([]);
  const [selectedPlatformId, setSelectedPlatformId] = useState<string>("");
  const [watchedDate, setWatchedDate] = useState<string>("");

  const token = useAuth();
  const accessToken = token ? token.accessToken : null;

  const headers: HeadersInit = {};
  if (accessToken) {
    headers.Authorization = `Bearer ${accessToken}`;
  }

  const authHeader = token ? `${token.tokenType} ${token.accessToken}` : "";
  const { toastData, captureBeforePoints, showToastAfter, closeToast } =
    usePointToast(authHeader);

  // movieDetailが更新されたときにrangeValueを設定
  useEffect(() => {
    if (movieDetail?.review?.rating) {
      onChange(movieDetail.review.rating.toString());
    }
  }, [movieDetail]);

  useEffect(() => {
    const fetchData = async () => {
      console.log("movieDetailのデータ取得中...");
      try {
        const [movieRes, platformRes] = await Promise.all([
          fetch(`/api/movies/${params.id}`, { method: "GET", headers }),
          fetch(`/api/platforms`, { method: "GET", headers }),
        ]);
        const movieDetail: MovieDetail = await movieRes.json();
        const platforms: Platform[] = await platformRes.json();

        console.log("movieDetailのデータ取得: 完了");
        setMovie(movieDetail);
        setPlatforms(platforms);
      } catch {
        console.log("データ取得: エラー");
        setMovie(undefined);
      }
    };

    fetchData();
  }, [params.id]);

  async function onSubmit(formData: FormData) {
    setIsLoading(true);
    try {
      // 視聴履歴（platformIdが選択されている場合のみ送信）
      const watchHistory = selectedPlatformId
        ? {
            platformId: Number(selectedPlatformId),
            watchedDate: watchedDate || undefined,
          }
        : undefined;

      const jsonData = {
        rating: formData.get("rating") || undefined,
        comment: formData.get("comment") || undefined,
        watchHistory,
      };

      const before = await captureBeforePoints();
      const response = await fetch(`/api/movies/${params.id}/reviews`, {
        method: "POST",
        headers,
        body: JSON.stringify(jsonData),
      });
      const resultCode: number = response.status;
      console.log("レビュー登録完了: %o", resultCode);

      if (resultCode === 201) {
        await showToastAfter(before);
        setTimeout(() => {
          router.push(`/movie/${params.id}?updated=true`);
        }, 2000);
      }
    } catch (error) {
      console.log(error);
    } finally {
      setIsLoading(false);
    }
  }

  if (!movieDetail) {
    return <div></div>;
  }

  // 今日の日付（視聴日の上限）
  const today = new Date().toISOString().split("T")[0];

  return (
    <div className="container-fluid pb-4">
      <PointToast data={toastData} onClose={closeToast} />
      <h3 className="text-center mb-4">Create Movie Review</h3>

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
                  sumReview={
                    movieDetail.voteCount
                      ? movieDetail.voteCount.toString()
                      : "0"
                  }
                />
              </div>
            </div>

            <div className="col-md-9">
              <div className="card-body text-light">
                {/* タイトル */}
                <h5 className="card-title">{movieDetail.title}</h5>
                {/* 自身の評価 */}
                {movieDetail.review?.rating && (
                  <div className="card-text">
                    <StarRating rating={movieDetail.review.rating} size={20} />
                  </div>
                )}
                {/* ジャンル */}
                {movieDetail.genres?.length !== 0 && (
                  <div className="card-text d-grid gap-2 d-md-block">
                    {movieDetail.genres?.map((g: Genre, i: number) => {
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
                    <div className="bg-transparent badge border border-info rounded-pill">
                      任意
                    </div>
                    {"　"}
                    プラットフォーム
                  </label>
                  <div className="col-sm-8 d-flex align-items-center">
                    <div className="px-2 row">
                      <div className="form-check form-check-inline col-md-3">
                        <input
                          className="form-check-input"
                          type="radio"
                          name="platformId"
                          id="platform-none"
                          value=""
                          checked={selectedPlatformId === ""}
                          onChange={() => setSelectedPlatformId("")}
                        />
                        <label
                          className="form-check-label"
                          htmlFor="platform-none"
                        >
                          登録しない
                        </label>
                      </div>
                      {platforms.map((p: Platform) => (
                        <div
                          key={p.code}
                          className="form-check form-check-inline col-md-3"
                        >
                          <input
                            className="form-check-input"
                            type="radio"
                            name="platformId"
                            id={p.code}
                            value={p.id}
                            checked={selectedPlatformId === p.id.toString()}
                            onChange={() =>
                              setSelectedPlatformId(p.id.toString())
                            }
                          />
                          <label className="form-check-label" htmlFor={p.code}>
                            {p.name}
                          </label>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>

                {/* 鑑賞日 */}
                {selectedPlatformId && (
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
                        value={watchedDate}
                        max={today}
                        onChange={(e) => setWatchedDate(e.target.value)}
                      />
                    </div>
                  </div>
                )}

                <div className="h4 pb-2 mb-4 text-success border-bottom border-success mt-4">
                  Review
                </div>

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
                      name="comment"
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
                    type="submit"
                    disabled={isLoading}
                    className="btn btn-outline-success"
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
