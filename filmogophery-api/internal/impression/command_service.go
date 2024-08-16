package impression

type (
	CommandService struct {
		ImpressionRepo ICommandRepository
	}
)

func NewCommandService(impressionRepo ICommandRepository) *CommandService {
	return &CommandService{
		ImpressionRepo: impressionRepo,
	}
}
