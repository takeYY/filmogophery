package mock

import "filmogophery/internal/app/types"

var (
	Action = types.Genre{
		Code: "action",
		Name: "アクション",
	}
	Adventure = types.Genre{
		Code: "adventure",
		Name: "アドベンチャー",
	}
	Comedy = types.Genre{
		Code: "comedy",
		Name: "コメディ",
	}
	Crime = types.Genre{
		Code: "crime",
		Name: "クライム",
	}
	Horror = types.Genre{
		Code: "horror",
		Name: "ホラー",
	}
	Romance = types.Genre{
		Code: "romance",
		Name: "ロマンス",
	}
	SF = types.Genre{
		Code: "sf",
		Name: "SF",
	}
	War = types.Genre{
		Code: "war",
		Name: "戦争",
	}
)
