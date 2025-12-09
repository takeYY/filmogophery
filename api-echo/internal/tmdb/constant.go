package tmdb

import "strconv"

func GetGenreName(tmdbGenreID *int) string {
	genres := map[string]string{
		"12":    "アドベンチャー",
		"14":    "ファンタジー",
		"16":    "アニメーション",
		"18":    "ドラマ",
		"27":    "ホラー",
		"28":    "アクション",
		"35":    "コメディ",
		"36":    "ヒストリー",
		"37":    "西部劇",
		"53":    "スリラー",
		"80":    "クライム",
		"99":    "ドキュメンタリー",
		"878":   "SF",
		"9648":  "ミステリー",
		"10402": "ミュージカル",
		"10749": "ロマンス",
		"10751": "ファミリー",
		"10752": "戦争",
		"10770": "TV",
	}
	key := strconv.Itoa(*tmdbGenreID)
	name, ok := genres[key]
	if !ok {
		return "不明"
	}

	return name
}
