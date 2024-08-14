import StarRating from "@/app/ui/Rating";
import { MovieDetail, WatchRecord } from "@/interface/movie";
import Image from "next/image";
import Link from "next/link";

export default function Page({ params }: { params: { id: string } }) {
  const movie = fetchMovie(params.id);
  return (
    <div className="container pb-4">
      <h3 className="text-center mb-4">Movie Detail</h3>
      <div
        className={`card mb-3 bg-dark ${
          movie.impression.status === "鑑賞済み" ? "border-success" : ""
        }`}
      >
        <div className="row g-0">
          <div className="col-md-3">
            {/* ポスター */}
            <Image
              src={movie.posterURL}
              className="img-fluid rounded-start"
              alt="ポスター画像"
              width={350}
              height={350}
            />

            {/* 一般の評価 */}
            <div className="justify-content-center">
              <StarRating
                rating={movie.vote_average}
                size={20}
                starColor={"#0dcaf0"}
                sumReview={movie.vote_count.toString()}
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
              {movie.genres.length && (
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
              <p className="card-text">
                公開日：{movie.release_date.toLocaleDateString()}
              </p>
              {/* 上映時間 */}
              <p className="card-text">上映時間：{movie.run_time}分</p>
              {/* 概要 */}
              <p className="card-text">{movie.overview}</p>
              {/* 感想 */}
              <div className="p-3 bg-success bg-opacity-10 border border-success border-start-0 border-end-0">
                {movie.impression.note}
              </div>
              {/* */}
            </div>
            {/* 視聴履歴 */}
            <div className="card-footer border-success text-light">
              <div>視聴履歴</div>
              {!movie.watchRecords.length && <div>なし</div>}

              {movie.watchRecords.length && (
                <dl className="row">
                  {movie.watchRecords.map((r: WatchRecord, i: number) => {
                    return (
                      <div key={i}>
                        <dt className="col-md-1 bg-transparent badge border border-primary rounded-pill">
                          {`${calcDiffDate(r.watch_date)}日前`}
                        </dt>
                        <dd className="col-md-10">
                          <dl className="row">
                            <dt className="col-md-4">
                              {r.watch_date.toLocaleDateString()}
                            </dt>
                            <dd className="col-md-8">{r.watch_media}</dd>
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
