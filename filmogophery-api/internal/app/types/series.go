package types

import "filmogophery/internal/pkg/gen/model"

type (
	Series struct {
		Name      string  `json:"name"`
		PosterURL *string `json:"posterURL"`
	}
)

func NewSeriesByModel(series *model.Series) *Series {
	if series == nil {
		return nil
	}
	return &Series{
		Name:      series.Name,
		PosterURL: series.PosterURL,
	}
}
