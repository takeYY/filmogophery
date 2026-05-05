package types

import "filmogophery/internal/pkg/gen/model"

type (
	Platform struct {
		ID   int32  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	}
)

func NewPlatformsByModel(platforms []*model.Platforms) []Platform {
	result := make([]Platform, 0, len(platforms))
	for _, p := range platforms {
		result = append(result, Platform{
			ID:   p.ID,
			Code: p.Code,
			Name: p.Name,
		})
	}

	return result
}

func NewPlatformByModel(platform model.Platforms) Platform {
	return Platform{
		ID:   platform.ID,
		Code: platform.Code,
		Name: platform.Name,
	}
}
