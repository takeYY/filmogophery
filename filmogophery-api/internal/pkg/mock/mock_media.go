package mock

import "filmogophery/internal/app/types"

var (
	PrimeVide  = types.Media{Code: "prime_video", Name: "Prime Video"}
	Netflix    = types.Media{Code: "netflix", Name: "Netflix"}
	UNext      = types.Media{Code: "u_next", Name: "U-NEXT"}
	DisneyPlus = types.Media{Code: "disney_plus", Name: "Disney+"}
	YouTube    = types.Media{Code: "youtube", Name: "YouTube"}
	AppleTV    = types.Media{Code: "apple_tv", Name: "Apple TV+"}
	Hulu       = types.Media{Code: "hulu", Name: "Hulu"}
	DAnime     = types.Media{Code: "d_anime", Name: "dアニメ"}
	Telasa     = types.Media{Code: "telasa", Name: "TELASA"}
	Cinema     = types.Media{Code: "cinema", Name: "映画館"}
	Unknown    = types.Media{Code: "unknown", Name: "不明"}

	MockedMedia = []types.Media{
		PrimeVide,
		Netflix,
		UNext,
		DisneyPlus,
		YouTube,
		AppleTV,
		Hulu,
		DAnime,
		Telasa,
		Cinema,
		Unknown,
	}
)
