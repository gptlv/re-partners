package api

import "github.com/gptlv/re-partners/packs/internal/app"

type Pack struct {
	Size  int64 `json:"size"`
	Count int64 `json:"count"`
}

type CalculateJSONRequest struct {
	Amount int64 `json:"amount"`
}

type CalculateJSONResponse struct {
	Amount int64  `json:"amount"`
	Packs  []Pack `json:"packs"`
}

func toAPIPacks(packs []app.Pack) []Pack {
	result := make([]Pack, 0, len(packs))
	for _, p := range packs {
		result = append(result, Pack{
			Size:  p.Size,
			Count: p.Count,
		})
	}
	return result
}
