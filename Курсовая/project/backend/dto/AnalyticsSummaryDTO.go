package dto

type PopularServices struct {
	Name    string `json:"name"`
	Count   int64  `json:"count"`
	Revenue string `json:"revenue"`
}

type DoctorLoad struct {
	Name       string `json:"name"`
	VisitCount int64  `json:"visitCount"`
	LoadStatus string `json:"loadStatus"`
}

type AnalyticsSummaryDTO struct {
	MonthlyVisits   int64             `json:"monthlyVisits"`
	TotalRevenue    string            `json:"totalRevenue"`
	RevenueChange   string            `json:"revenueChange"`
	ActiveClients   int64             `json:"activeClients"`
	AvgCheck        int64             `json:"avgCheck"`
	PopularServices []PopularServices `json:"popularServices"`
	DoctorLoad      []DoctorLoad      `json:"doctorLoad"`
}
