package model

type Price struct {
	ID               string `json:"id"`
	PriceType        string `json:"priceType"`
	PriceNum         string `json:"priceNum"`
	PriceDen         string `json:"priceDen"`
	PriceDec         string `json:"priceDec"`
	PriceAmerican    string `json:"priceAmerican"`
	HandicapValueDec string `json:"handicapValueDec"`
	RawHandicapValue string `json:"rawHandicapValue"`
	IsActive         string `json:"isActive"`
	DisplayOrder     string `json:"displayOrder"`
}

type Outcome struct {
	ID                       string `json:"id"`
	MarketID                 string `json:"marketId"`
	Name                     string `json:"name"`
	OutcomeMeaningMajorCode  string `json:"outcomeMeaningMajorCode"`
	OutcomeMeaningMinorCode  string `json:"outcomeMeaningMinorCode"`
	DisplayOrder             string `json:"displayOrder"`
	OutcomeStatusCode        string `json:"outcomeStatusCode"`
	IsActive                 string `json:"isActive"`
	IsDisplayed              string `json:"isDisplayed"`
	SiteChannels             string `json:"siteChannels"`
	FixedOddsAvail           string `json:"fixedOddsAvail"`
	LiveServChannels         string `json:"liveServChannels"`
	LiveServChildrenChannels string `json:"liveServChildrenChannels"`
	IsAvailable              string `json:"isAvailable"`
	CashoutAvail             string `json:"cashoutAvail"`
	Children                 []struct {
		Price Price `json:"price"`
	} `json:"children"`
}

type Market struct {
	ID                       string `json:"id"`
	EventID                  string `json:"eventId"`
	TemplateMarketID         string `json:"templateMarketId"`
	TemplateMarketName       string `json:"templateMarketName"`
	MarketMeaningMajorCode   string `json:"marketMeaningMajorCode"`
	MarketMeaningMinorCode   string `json:"marketMeaningMinorCode"`
	Name                     string `json:"name"`
	IsLpAvailable            string `json:"isLpAvailable"`
	RawHandicapValue         string `json:"rawHandicapValue"`
	DisplayOrder             string `json:"displayOrder"`
	MarketStatusCode         string `json:"marketStatusCode"`
	IsActive                 string `json:"isActive"`
	IsDisplayed              string `json:"isDisplayed"`
	FixedOddsAvail           string `json:"fixedOddsAvail"`
	SiteChannels             string `json:"siteChannels"`
	LiveServChannels         string `json:"liveServChannels"`
	LiveServChildrenChannels string `json:"liveServChildrenChannels"`
	PriceTypeCodes           string `json:"priceTypeCodes"`
	IsBettable               string `json:"isBettable"`
	IsAvailable              string `json:"isAvailable"`
	MaxAccumulators          string `json:"maxAccumulators"`
	MinAccumulators          string `json:"minAccumulators"`
	CashoutAvail             string `json:"cashoutAvail"`
	TermsWithBet             string `json:"termsWithBet"`
	Children                 []struct {
		Outcome Outcome `json:"outcome"`
	} `json:"children"`
}

type Event struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	EventStatusCode          string `json:"eventStatusCode"`
	IsActive                 string `json:"isActive"`
	IsDisplayed              string `json:"isDisplayed"`
	DisplayOrder             string `json:"displayOrder"`
	SiteChannels             string `json:"siteChannels"`
	EventSortCode            string `json:"eventSortCode"`
	StartTime                string `json:"startTime"`
	SuspendAtTime            string `json:"suspendAtTime"`
	RawIsOffCode             string `json:"rawIsOffCode"`
	ClassID                  string `json:"classId"`
	TypeID                   string `json:"typeId"`
	SportID                  string `json:"sportId"`
	LiveServChannels         string `json:"liveServChannels"`
	LiveServChildrenChannels string `json:"liveServChildrenChannels"`
	CategoryID               string `json:"categoryId"`
	CategoryCode             string `json:"categoryCode"`
	CategoryName             string `json:"categoryName"`
	CategoryDisplayOrder     string `json:"categoryDisplayOrder"`
	ClassName                string `json:"className"`
	ClassDisplayOrder        string `json:"classDisplayOrder"`
	ClassSortCode            string `json:"classSortCode"`
	TypeName                 string `json:"typeName"`
	TypeDisplayOrder         string `json:"typeDisplayOrder"`
	TypeFlagCodes            string `json:"typeFlagCodes"`
	IsOpenEvent              string `json:"isOpenEvent"`
	IsAvailable              string `json:"isAvailable"`
	FixedOddsAvail           string `json:"fixedOddsAvail"`
	CashoutAvail             string `json:"cashoutAvail"`
	IsBettable               string `json:"isBettable"`
	Children                 []struct {
		Market Market `json:"market"`
	} `json:"children"`
}

type EventsRoot struct {
	SSResponse struct {
		Children []struct {
			Event Event `json:"event"`
		} `json:"children"`
	} `json:"SSResponse"`
}
