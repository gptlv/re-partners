package api

import "github.com/gptlv/re-partners/packs/internal/app"

type Pack struct {
	Size  int64 `json:"size"`
	Count int64 `json:"count"`
}

type PackSizeResponse struct {
	ID   int64 `json:"id"`
	Size int64 `json:"size"`
}

type PackSizesResponse struct {
	Packs []PackSizeResponse `json:"packs"`
}

type CalculateRequest struct {
	Amount int64 `json:"amount"`
}

type CalculateResponse struct {
	Amount int64  `json:"amount"`
	Packs  []Pack `json:"packs"`
}

type CreateSizeRequest struct {
	Size int64 `json:"size"`
}

type CreateSizeResponse struct {
	ID   int64 `json:"id"`
	Size int64 `json:"size"`
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

func toPackSizeResponses(sizes []app.PackSize) []PackSizeResponse {
	result := make([]PackSizeResponse, 0, len(sizes))
	for _, s := range sizes {
		result = append(result, PackSizeResponse{
			ID:   s.ID,
			Size: s.Size,
		})
	}
	return result
}
