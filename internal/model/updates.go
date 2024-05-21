package model

type EventUpdate struct {
	Names        map[string]string `json:"names"`
	StartTimeXls map[string]string `json:"start_time_xls"`
	Status       string            `json:"status"`
	Displayed    string            `json:"displayed"`
	ResultConf   string            `json:"result_conf"`
	StartTime    string            `json:"start_time"`
	SuspendAt    string            `json:"suspend_at"`
	IsOff        string            `json:"is_off"`
	Started      string            `json:"started"`
	RaceStage    string            `json:"race_stage"`
	Disporder    int               `json:"disporder"`
}

type Collection struct {
	CollectionID string `json:"collection_id"`
}

type MarketUpdate struct {
	Names                map[string]string `json:"names"`
	GroupNames           map[string]string `json:"group_names"`
	HcapValues           map[string]string `json:"hcap_values"`
	EvOcGrpID            string            `json:"ev_oc_grp_id"`
	MktDispCode          string            `json:"mkt_disp_code"`
	MktDispLayoutColumns string            `json:"mkt_disp_layout_columns"`
	MktDispLayoutOrder   string            `json:"mkt_disp_layout_order"`
	MktType              string            `json:"mkt_type"`
	MktSort              string            `json:"mkt_sort"`
	MktGrpFlags          string            `json:"mkt_grp_flags"`
	Status               string            `json:"status"`
	Displayed            string            `json:"displayed"`
	BirIndex             string            `json:"bir_index"`
	RawHcap              string            `json:"raw_hcap"`
	EwAvail              string            `json:"ew_avail"`
	EwPlaces             string            `json:"ew_places"`
	EwFacNum             string            `json:"ew_fac_num"`
	EwFacDen             string            `json:"ew_fac_den"`
	BetInRun             string            `json:"bet_in_run"`
	LpAvail              string            `json:"lp_avail"`
	SpAvail              string            `json:"sp_avail"`
	FcAvail              string            `json:"fc_avail"`
	TcAvail              string            `json:"tc_avail"`
	MmCollID             string            `json:"mm_coll_id"`
	SuspendAt            string            `json:"suspend_at"`
	Collections          []Collection      `json:"collections"`
	EvID                 int               `json:"ev_id"`
	Disporder            int               `json:"disporder"`
}

type SelectionUpdate struct {
	Names     map[string]string `json:"names"`
	Status    string            `json:"status"`
	Settled   string            `json:"settled"`
	Result    string            `json:"result"`
	Displayed string            `json:"displayed"`
	RunnerNum string            `json:"runner_num"`
	FbResult  string            `json:"fb_result"`
	LpNum     string            `json:"lp_num"`
	LpDen     string            `json:"lp_den"`
	CsHome    string            `json:"cs_home"`
	CsAway    string            `json:"cs_away"`
	Flags     string            `json:"flags"`
	UniqueID  string            `json:"unique_id"`
	EvMktID   int               `json:"ev_mkt_id"`
	Disporder int               `json:"disporder"`
}

type PriceUpdate struct {
	LpNum string `json:"lp_num"`
	LpDen string `json:"lp_den"`
}
