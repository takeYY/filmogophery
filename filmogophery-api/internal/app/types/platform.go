package types

import "filmogophery/internal/pkg/gen/model"

type (
	Platform struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}
)

func NewPlatformsByModel(platforms []*model.Platforms) []Platform {
	result := make([]Platform, 0, len(platforms))
	for _, p := range platforms {
		result = append(result, Platform{
			Code: p.Code,
			Name: p.Name,
		})
	}

	return result
}
