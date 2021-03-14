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
	summaries []model.TossSummary,
) []GetSuggestHandResponse {
	var resp []GetSuggestHandResponse
	for i := range summaries {
		summ := summaries[i]
		hand := make([]string, len(summ.Kept))
		for i, c := range summ.Kept {
			hand[i] = c.String()
		}

		toss := make([]string, len(summ.Tossed))
		for i, c := range summ.Tossed {
			toss[i] = c.String()
		}

		resp = append(resp, GetSuggestHandResponse{
			Hand: hand,
			Toss: toss,
			HandPts: PointStats{
				Min:    summ.HandStats.Min(),
				Avg:    summ.HandStats.Avg(),
				Median: summ.HandStats.Median(),
				Max:    summ.HandStats.Max(),
			},
			CribPts: PointStats{
				Min:    summ.CribStats.Min(),
				Avg:    summ.CribStats.Avg(),
				Median: summ.CribStats.Median(),
				Max:    summ.CribStats.Max(),
			},
		})
	}
	return resp
}
