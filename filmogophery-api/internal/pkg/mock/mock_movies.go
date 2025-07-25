package mock

import "filmogophery/internal/app/types"

var (
	posterURL001 = "/iK2YBfD7DdaNXZALQhhT9uTN9Rc.jpg"
	posterURL002 = "/oCwo7ALD3LftLqy0Oj6U669u4fU.jpg"
	posterURL004 = "/7fmBDyHsLTyhfRSZk19wrY23zDg.jpg"
	posterURL006 = "/bOGsmFU5Lc9U56Tq01gFm6wbZnb.jpg"
	posterURL007 = "/g9oEopuLdvVZNZw54vPC2fe6tNv.jpg"
	posterURL008 = "/fzdmsB3rcxu6wX5nZqRE6H92VAn.jpg"
	MockedMovies = []types.Movie{
		{
			ID:    218,
			Title: "ターミネーター",
			Overview: `アメリカのとある街、深夜突如奇怪な放電と共に屈強な肉体をもった男が現れる。
同じくして放電の中からもう一人の男カイル・リースが現れる。
屈強な肉体を持った男はモラルや常識もない。
あるのはただ１つの目的アメリカ人女性サラ・コナーという名の人物の殺害だった。
電話帳名簿から「サラ・コナー」の名を持つ女性をかたっぱしから銃殺していく男。
その頃カイルは目的のサラ・コナーと接触し間一髪で彼女を救う。
カイルはサラに、サラを狙っているのは近未来から送られた人類殺戮ロボット「ターミネーター」であり、
未来ではロボットの反乱による機械対人類の最終戦争が起こっている事、
そしてサラは人類軍の希望のリーダー、ジョン・コナーの母親である事を告げる。
サラを連れカイルはターミネーターからの逃亡を開始する。`,
			ReleaseDate: "1985-05-04",
			RunTime:     108,
			PosterURL:   &posterURL001,
			TmdbID:      218,
			Genres:      []types.Genre{Action, Horror, SF},
		},
		{
			ID:    280,
			Title: "ターミネーター2",
			Overview: `未来からの抹殺兵器ターミネーターを破壊し、近未来で恐ろしい戦争が起こる事を知ってしまったサラ・コナー。
カイルとの子供ジョンは母親から常にその戦争の話や戦いへの備えの話を聞かされていた。
サラは周囲から変人扱いされ精神病院へ収容されジョンは親戚の家で暮らしていた。
ある日ジョンの前に執拗にジョンを狙う不審な警官が現る。
軌道を逸した警官の行動は明らかにジョンを殺害しようとしていた。
殺されるその寸前、見知らぬ屈強な男が現れジョンを救う。
彼は自らをターミネーターでありジョンを守るべく再プログラムされ未来から送り込まれたと告げる。
ジョンは病院の母親を連れ出し、最終戦争を起こす原因であるコンピューターシステム「スカイネット」を破壊するため、
ターミネーターと共にサイバーダイン社へ向かうが・・・`,
			ReleaseDate: "1991-08-24",
			RunTime:     137,
			PosterURL:   &posterURL002,
			TmdbID:      280,
			Genres:      []types.Genre{},
		},
		{
			ID:    296,
			Title: "ターミネーター3",
			Overview: `スカイネットを破壊、予見されていた最終戦争の日も過ぎたが、
ジョンの心には母親から刷り込まれた戦争の話が無くなったとは思えず不安な毎日を過ごしていた。
そんな中、未来より新たなターミネーターがジョンの仲間となる兵士達の抹殺のために送られてくる。
新型ターミネーターT-Xは、ジョンを発見したことで、ジョン自身の殺害プログラムを優先し彼を追うが、
たまたまいっしょにいた女性と共に別のターミネーターT-850に救われる。
T-850は未来のジョンの妻ケイトによりリプログラムされ送られたのだった。
ジョンとケイトは最終戦争開始の日である審判の日は回避することはできず、
スカイネット自体も破壊を免れており起動を待ち軍事基地に現存している事を知らされる。
スカイネット起動時間が迫るなか、ジョンとケイトは起動を阻止すべくT-850の指示のもと、アメリカ軍基地へ向かうが・・・`,
			ReleaseDate: "2003-07-12",
			RunTime:     110,
			PosterURL:   nil,
			TmdbID:      296,
			Genres:      []types.Genre{},
		},
		{
			ID:    534,
			Title: "ターミネーター4",
			Overview: `“審判の日”から10年後の2018年。
人類軍の指導者となり、機械軍と戦うことを幼いころから運命づけられてきたジョン・コナー。
今や30代となった彼は、人類滅亡をもくろむスカイネットの猛攻が開始されようとする中、
ついに人類軍のリーダーとして立ち上がることになる。`,
			ReleaseDate: "2009-06-05",
			RunTime:     114,
			PosterURL:   &posterURL004,
			TmdbID:      534,
			Genres:      []types.Genre{Action, Horror, SF, War},
		},
		{
			ID:    105,
			Title: "バック・トゥ・ザ・フューチャー",
			Overview: `スティーブン・スピルバーグとロバート・ゼメキスが贈るSFアドベンチャーシリーズ第1弾。
高校生のマーティは、科学者・ドクの発明したタイムマシン・デロリアンで過去にタイムスリップしてしまう。`,
			ReleaseDate: "1985-12-07",
			RunTime:     116,
			PosterURL:   nil,
			TmdbID:      105,
			Genres:      []types.Genre{Action, Adventure, Comedy, Romance, SF},
		},
		{
			ID:    165,
			Title: "バック・トゥ・ザ・フューチャー PART2",
			Overview: `スティーブン・スピルバーグとロバート・ゼメキスが贈るSFアドベンチャーシリーズ第2弾。
現代に戻って来たマーティは、2015年から帰って来たドクに連れられ今度は未来へタイムスリップすることに。`,
			ReleaseDate: "1989-12-09",
			RunTime:     108,
			PosterURL:   &posterURL006,
			TmdbID:      165,
			Genres:      []types.Genre{},
		},
		{
			ID:    19995,
			Title: "アバター",
			Overview: `西暦2154年。人類は惑星ポリフェマスの最大衛星パンドラに鉱物採掘基地を開いている。
この星は熱帯雨林のような未開のジャングル覆われていて獰猛な動物と”ナヴィ”という先住種族が暮らしており、
森の奥には地球のエネルギー問題解決の鍵となる希少鉱物の鉱床がある。
この星の大気は人間に適さないので屋外活動にはマスクを着用する必要があり、
ナヴィと意思疎通し交渉するために人間とナヴィの遺伝子を組み合わせて人間が作りあげた”アバター”が用いられた。`,
			ReleaseDate: "2010-08-26",
			RunTime:     162,
			PosterURL:   &posterURL007,
			TmdbID:      19995,
			Genres:      []types.Genre{SF},
		},
		{
			ID:    76600,
			Title: "アバター：ウェイ・オブ・ウォーター",
			Overview: `第1作目から約10年後の惑星パンドラでのジェイクとネイティリの子供たちからなる家族の物語。
一家は神聖なる森を追われ海の部族に助けを求めるが、その楽園のような海辺の世界にも人類の侵略の手が迫っていた・・・。`,
			ReleaseDate: "2022-12-16",
			RunTime:     192,
			PosterURL:   &posterURL008,
			TmdbID:      76600,
			Genres:      []types.Genre{SF},
		},
	}
	movie218  = MockedMovies[0]
	rating218 = float32(4.3)
	note218   = `ターミネーターの元祖という感じで、恐ろしさと希望が織り成す圧巻の作品。今観るとVFXのぎこちなさが目立つが、それが逆に怖さを演出していて好印象。`

	movie280                = MockedMovies[1]
	MockedMovieDetailMapper = map[int32]types.MovieDetail{
		218: {
			VoteAverage: float32(77) / 10 / 2,
			VoteCount:   12951,
			Series:      nil,
			Impression: &types.Impression{
				ID:     218,
				Status: "",
				Rating: &rating218,
				Note:   &note218,
				Records: []types.Record{
					{
						WatchDate:  "2016-12-25",
						WatchMedia: "不明",
					},
					{
						WatchDate:  "2022-10-24",
						WatchMedia: "Prime Video",
					},
					{
						WatchDate:  "2024-08-01",
						WatchMedia: "U-NEXT",
					},
				},
			},
			Movie: movie218,
		},
		280: {
			VoteAverage: float32(81) / 10 / 2,
			VoteCount:   12699,
			Series:      nil,
			Impression: &types.Impression{
				ID:      280,
				Status:  "未鑑賞",
				Rating:  nil,
				Note:    nil,
				Records: []types.Record{},
			},
			Movie: movie280,
		},
	}
)
