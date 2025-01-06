package impression

/*
type InMemoryRepository struct {
	movieImpressions []*model.MovieImpression
	mu               sync.RWMutex
}

func NewInMemoryRepository(testData []*model.MovieImpression) *impression.IQueryRepository {
	var inMemoryQueryRepo impression.IQueryRepository = &InMemoryRepository{
		movieImpressions: testData,
	}
	return &inMemoryQueryRepo
}

func (ir *InMemoryRepository) Find(ctx context.Context) ([]*model.MovieImpression, error) {
	ir.mu.RLock()
	defer ir.mu.RUnlock()

	impressions := ir.movieImpressions

	return impressions, nil
}
*/
