import { MovieDetail, Impression, WatchRecord } from "@/interface/movie";

export const impressionData: Map<number, Impression> = new Map([
  [
    218,
    {
      id: 218,
      status: "",
      rating: 4.3,
      note: `ターミネーターの元祖という感じで、恐ろしさと希望が織り成す圧巻の作品。
  今観るとVFXのぎこちなさが目立つが、それが逆に怖さを演出していて好印象。`,
    },
  ],
  [
    280,
    {
      id: 280,
      status: "未鑑賞",
      rating: null,
      note: null,
    },
  ],
]);

const watchRecordData: Map<number, WatchRecord[]> = new Map([
  [
    218,
    [
      {
        watchDate: "2016-12-25",
        watchMedia: "不明",
      },
      {
        watchDate: "2022-10-24",
        watchMedia: "Prime Video",
      },
      {
        watchDate: "2024-08-01",
        watchMedia: "U-NEXT",
      },
    ],
  ],
]);

export const movieDetailData: Map<string, MovieDetail> = new Map([
  [
    "218",
    {
      id: 218,
      title: "ターミネーター",
      overview: `アメリカのとある街、深夜突如奇怪な放電と共に屈強な肉体をもった男が現れる。
        同じくして放電の中からもう一人の男カイル・リースが現れる。
        屈強な肉体を持った男はモラルや常識もない。
        あるのはただ１つの目的アメリカ人女性サラ・コナーという名の人物の殺害だった。
        電話帳名簿から「サラ・コナー」の名を持つ女性をかたっぱしから銃殺していく男。
        その頃カイルは目的のサラ・コナーと接触し間一髪で彼女を救う。
        カイルはサラに、サラを狙っているのは近未来から送られた人類殺戮ロボット「ターミネーター」であり、
        未来ではロボットの反乱による機械対人類の最終戦争が起こっている事、
        そしてサラは人類軍の希望のリーダー、ジョン・コナーの母親である事を告げる。
        サラを連れカイルはターミネーターからの逃亡を開始する。`,
      releaseDate: "1985-05-04",
      runTime: 108,
      genres: ["アクション", "ホラー", "SF"],
      posterURL: "/iK2YBfD7DdaNXZALQhhT9uTN9Rc.jpg",
      voteAverage: 77 / 10 / 2,
      voteCount: 12951,
      series: null,
      impression: impressionData.get(218) ?? null,
      watchRecords: watchRecordData.get(218) ?? [],
    },
  ],
  [
    "280",
    {
      id: 280,
      title: "ターミネーター2",
      overview: `未来からの抹殺兵器ターミネーターを破壊し、近未来で恐ろしい戦争が起こる事を知ってしまったサラ・コナー。
    カイルとの子供ジョンは母親から常にその戦争の話や戦いへの備えの話を聞かされていた。
    サラは周囲から変人扱いされ精神病院へ収容されジョンは親戚の家で暮らしていた。
    ある日ジョンの前に執拗にジョンを狙う不審な警官が現る。
    軌道を逸した警官の行動は明らかにジョンを殺害しようとしていた。
    殺されるその寸前、見知らぬ屈強な男が現れジョンを救う。
    彼は自らをターミネーターでありジョンを守るべく再プログラムされ未来から送り込まれたと告げる。
    ジョンは病院の母親を連れ出し、最終戦争を起こす原因であるコンピューターシステム「スカイネット」を破壊するため、
    ターミネーターと共にサイバーダイン社へ向かうが・・・`,
      releaseDate: "1991-08-24",
      runTime: 137,
      genres: [],
      posterURL: "/oCwo7ALD3LftLqy0Oj6U669u4fU.jpg",
      voteAverage: 81 / 10 / 2,
      voteCount: 12699,
      series: null,
      impression: impressionData.get(280) ?? null,
      watchRecords: [],
    },
  ],
]);
