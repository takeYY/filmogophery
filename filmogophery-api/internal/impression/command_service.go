package impression

import "filmogophery/internal/app/repositories"

type (
	CommandService struct {
		ImpressionRepo repositories.IImpressionRepository
	}
)

func NewCommandService(impressionRepo repositories.IImpressionRepository) *CommandService {
	return &CommandService{
		ImpressionRepo: impressionRepo,
	}
}
