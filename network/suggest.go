package network

import "github.com/joshprzybyszewski/cribbage/model"

type PointStats struct {
	Min    int     `json:"min"`
	Median float64 `json:"median"`
	Avg    float64 `json:"avg"`
	Max    int     `json:"max"`
}

type GetSuggestHandResponse struct {
	Hand    []string   `json:"hand"`
	Toss    []string   `json:"toss"`
	HandPts PointStats `json:"handPts"`
	CribPts PointStats `json:"cribPts"`
}

func ConvertToGetSuggestHandResponse(
	sums []model.TossSummary,
) []GetSuggestHandResponse {
	var resp []GetSuggestHandResponse
	for _, sum := range sums {
		hand := make([]string, len(sum.Kept))
		for i, c := range sum.Kept {
			hand[i] = c.String()
		}

		toss := make([]string, len(sum.Tossed))
		for i, c := range sum.Tossed {
			toss[i] = c.String()
		}

		resp = append(resp, GetSuggestHandResponse{
			Hand: hand,
			Toss: toss,
			HandPts: PointStats{
				Min:    sum.HandStats.Min(),
				Avg:    sum.HandStats.Avg(),
				Median: sum.HandStats.Median(),
				Max:    sum.HandStats.Max(),
			},
			CribPts: PointStats{
				Min:    sum.CribStats.Min(),
				Avg:    sum.CribStats.Avg(),
				Median: sum.CribStats.Median(),
				Max:    sum.CribStats.Max(),
			},
		})
	}
	return resp
}
