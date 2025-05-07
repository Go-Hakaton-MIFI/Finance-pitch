package schemas

type DateRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type AnalyticsFilter struct {
	Date DateRange `json:"date" validate:"required"`
}

type DynamicsByPeriodRequest struct {
	Date DateRange `json:"date" validate:"required"`
}

type DynamicsByPeriodResponse struct {
	Data []struct {
		Date  string  `json:"date"`
		Value float64 `json:"value"`
	} `json:"data"`
}

type CategoriesSummaryRequest struct {
	Date DateRange `json:"date" validate:"required"`
}

type CategoriesSummaryResponse struct {
	Data []struct {
		Category string  `json:"category"`
		Value    float64 `json:"value"`
	} `json:"data"`
}
