package services

type (
	IWatchlistService interface {
		// --- Create --- //

		// --- Read --- //

		// --- Update --- //

		// --- Delete --- //
	}

	watchlistService struct{}
)

func NewWatchlistService() IWatchlistService {
	return &watchlistService{}
}
