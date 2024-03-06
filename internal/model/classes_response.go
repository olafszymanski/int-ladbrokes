package model

type Class struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	ClassStatusCode         string `json:"classStatusCode"`
	IsActive                string `json:"isActive"`
	DisplayOrder            string `json:"displayOrder"`
	SiteChannels            string `json:"siteChannels"`
	ClassSortCode           string `json:"classSortCode"`
	CategoryID              string `json:"categoryId"`
	CategoryCode            string `json:"categoryCode"`
	CategoryName            string `json:"categoryName"`
	CategoryDisplayOrder    string `json:"categoryDisplayOrder"`
	HasOpenEvent            string `json:"hasOpenEvent"`
	HasNext24HourEvent      string `json:"hasNext24HourEvent"`
	HasLiveNowOrFutureEvent string `json:"hasLiveNowOrFutureEvent"`
}

type ClassesRoot struct {
	SSResponse struct {
		Children []struct {
			Class Class `json:"class"`
		} `json:"children"`
	} `json:"SSResponse"`
}
