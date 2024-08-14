"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import { Movie } from "@/interface/movie";

export default function Home() {
  const router = useRouter();
  const movies = fetchMovies();
  return (
    <main>
      <div className="container pb-4">
        <h3 className="text-center mb-4">Home</h3>

        {/* TODO: レイアウトが崩れているので、直すこと!! */}
        <div className="card-columns">
          {movies &&
            movies.map((obj: Movie, index: number) => {
              return (
                <button
                  key={index}
                  className="card bg-dark"
                  onClick={() => router.push(`/movie/${obj.id}`)}
                >
                  <div className="row no-gutters">
                    {/* ポスター */}
                    <div className="col-md-4">
                      <Image
                        src={obj.posterURL}
                        alt="ポスター画像"
                        className="img-fluid"
                        width={75}
                        height={75}
                      />
                    </div>
                    <div className="col-md-8">
                      <div className="card-body">
                        {/* タイトル */}
                        <h5 className="card-title text-light">{obj.title}</h5>
                        {/* 概要 */}
                        <p className="card-text text-light">{obj.overview}</p>
                        {/* ジャンル */}
                        {obj.genres &&
                          obj.genres.map((g: string, i: number) => {
                            return (
                              <p key={i} className="badge text-bg-secondary">
                                {g}
                              </p>
                            );
                          })}
                      </div>
                    </div>
                  </div>
                </button>
              );
            })}
        </div>
      </div>
    </main>
  );
}

function fetchMovies(): Movie[] | undefined {
  try {
    console.log("データ取得中...");
    const movies: Movie[] = [
      {
        id: 1,
        title: "ターミネーター",
        overview:
          "アメリカのとある街、深夜突如奇怪な放電と共に屈強な肉体をもった男が現れる。同じく...",
        release_date: "1985-05-04",
        run_time: 108,
        genres: ["アクション", "SF"],
        posterURL:
          "https://image.tmdb.org/t/p/w600_and_h900_bestv2/iK2YBfD7DdaNXZALQhhT9uTN9Rc.jpg",
      },
      {
        id: 2,
        title: "ターミネーター2",
        overview:
          "未来からの抹殺兵器ターミネーターを破壊し、近未来で恐ろしい戦争が起こる事を知って...",
        release_date: "1991-08-24",
        run_time: 137,
        genres: [],
        posterURL:
          "https://image.tmdb.org/t/p/w600_and_h900_bestv2/oCwo7ALD3LftLqy0Oj6U669u4fU.jpg",
      },
      {
        id: 3,
        title: "ターミネーター3",
        overview:
          "スカイネットを破壊、予見されていた最終戦争の日も過ぎたが、ジョンの心には母親から...",
        release_date: "2003-07-12",
        run_time: 110,
        genres: [],
        posterURL:
          "https://image.tmdb.org/t/p/w600_and_h900_bestv2/oCwo7ALD3LftLqy0Oj6U669u4fU.jpg",
      },
      {
        id: 4,
        title: "ターミネーター4",
        overview:
          "“審判の日”から10年後の2018年。人類軍の指導者となり、機械軍と戦うことを幼...",
        release_date: "2009-06-05",
        run_time: 114,
        genres: ["アクション", "SF", "戦争"],
        posterURL:
          "https://image.tmdb.org/t/p/w600_and_h900_bestv2/7fmBDyHsLTyhfRSZk19wrY23zDg.jpg",
      },
    ];
    // const response = await fetch(`http://localhost:8000/movies`);
    // const data = await response.json();
    console.log("データ取得: 完了!!!");
    return movies;
  } catch (error) {
    console.log("データ取得エラー");
  }
}
