package dto

type DailyRows struct {
	Date     string `json:"date"`
	Visits   string `json:"visits"`
	Services string `json:"services"`
	Meds     string `json:"meds"`
	Total    string `json:"total"`
}

type RevenueReportDTO struct {
	PeriodTotal   int64       `json:"periodTotal"`
	ServicesTotal int64       `json:"servicesTotal"`
	DailyRows     []DailyRows `json:"dailyRows"`
}
