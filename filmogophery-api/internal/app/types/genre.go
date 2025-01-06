package types

import "filmogophery/internal/pkg/gen/model"

type (
	Genre struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}
)

func NewGenresByModel(genres []*model.Genre) []Genre {
	result := make([]Genre, 0, len(genres))
	for _, g := range genres {
		result = append(result, Genre{
			Code: g.Code,
			Name: *g.Name,
		})
	}

	return result
}
