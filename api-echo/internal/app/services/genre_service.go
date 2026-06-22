package services

import (
	"context"

	"github.com/rs/zerolog"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/pkg/gen/model"
)

type (
	IGenreService interface {
		// --- Read --- //

		// 全てのジャンルを取得
		GetAllGenres(ctx context.Context) ([]*model.Genres, error)
	}

	genreService struct {
		genreRepo repositories.IGenreRepository
	}
)

func NewGenreService(
	genreRepo repositories.IGenreRepository,
) IGenreService {
	return &genreService{
		genreRepo,
	}
}

// 全てのジャンルを取得
func (s *genreService) GetAllGenres(ctx context.Context) ([]*model.Genres, error) {
	log := zerolog.Ctx(ctx)

	genres, err := s.genreRepo.FindAll(ctx)
	if err != nil {
		log.Error().Msgf("failed to fetch genres: %s", err.Error())
		return nil, responses.InternalServerError()
	}

	return genres, nil
}
