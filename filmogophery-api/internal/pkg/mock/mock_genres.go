package mock

import "filmogophery/internal/app/types"

var (
	Action      = types.Genre{Code: "action", Name: "アクション"}
	Adventure   = types.Genre{Code: "adventure", Name: "アドベンチャー"}
	Animation   = types.Genre{Code: "animation", Name: "アニメーション"}
	Comedy      = types.Genre{Code: "comedy", Name: "コメディ"}
	Crime       = types.Genre{Code: "crime", Name: "クライム"}
	Documentary = types.Genre{Code: "documentary", Name: "ドキュメンタリー"}
	Drama       = types.Genre{Code: "drama", Name: "ドラマ"}
	Family      = types.Genre{Code: "family", Name: "ファミリー"}
	Fantasy     = types.Genre{Code: "fantasy", Name: "ファンタジー"}
	History     = types.Genre{Code: "history", Name: "ヒストリー"}
	Horror      = types.Genre{Code: "horror", Name: "ホラー"}
	Musical     = types.Genre{Code: "musical", Name: "ミュージカル"}
	Mystery     = types.Genre{Code: "mystery", Name: "ミステリー"}
	Romance     = types.Genre{Code: "romance", Name: "ロマンス"}
	SF          = types.Genre{Code: "sf", Name: "SF"}
	TV          = types.Genre{Code: "tv", Name: "TV"}
	Thriller    = types.Genre{Code: "thriller", Name: "スリラー"}
	War         = types.Genre{Code: "war", Name: "戦争"}
	Western     = types.Genre{Code: "western", Name: "西部劇"}

	MockedGenres = []types.Genre{
		Action,
		Adventure,
		Animation,
		Comedy,
		Crime,
		Documentary,
		Drama,
		Family,
		Fantasy,
		History,
		Horror,
		Musical,
		Mystery,
		Romance,
		SF,
		TV,
		Thriller,
		War,
		Western,
	}
)
