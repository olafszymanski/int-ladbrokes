package model

type EventUpdate struct {
	Names        map[string]string `json:"names"`
	Status       string            `json:"status"`
	Displayed    string            `json:"displayed"`
	ResultConf   string            `json:"result_conf"`
	Disporder    int               `json:"disporder"`
	StartTime    string            `json:"start_time"`
	StartTimeXls map[string]string `json:"start_time_xls"`
	SuspendAt    string            `json:"suspend_at"`
	IsOff        string            `json:"is_off"`
	Started      string            `json:"started"`
	RaceStage    string            `json:"race_stage"`
}

type Collection struct {
	CollectionID string `json:"collection_id"`
}

type MarketUpdate struct {
	Names                map[string]string `json:"names"`
	GroupNames           map[string]string `json:"group_names"`
	EvOcGrpID            string            `json:"ev_oc_grp_id"`
	MktDispCode          string            `json:"mkt_disp_code"`
	MktDispLayoutColumns string            `json:"mkt_disp_layout_columns"`
	MktDispLayoutOrder   string            `json:"mkt_disp_layout_order"`
	MktType              string            `json:"mkt_type"`
	MktSort              string            `json:"mkt_sort"`
	MktGrpFlags          string            `json:"mkt_grp_flags"`
	EvID                 int               `json:"ev_id"`
	Status               string            `json:"status"`
	Displayed            string            `json:"displayed"`
	Disporder            int               `json:"disporder"`
	BirIndex             string            `json:"bir_index"`
	RawHcap              string            `json:"raw_hcap"`
	HcapValues           map[string]string `json:"hcap_values"`
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
}

type SelectionUpdate struct {
	Names     map[string]string `json:"names"`
	EvMktID   int               `json:"ev_mkt_id"`
	Status    string            `json:"status"`
	Settled   string            `json:"settled"`
	Result    string            `json:"result"`
	Displayed string            `json:"displayed"`
	Disporder int               `json:"disporder"`
	RunnerNum string            `json:"runner_num"`
	FbResult  string            `json:"fb_result"`
	LpNum     string            `json:"lp_num"`
	LpDen     string            `json:"lp_den"`
	CsHome    string            `json:"cs_home"`
	CsAway    string            `json:"cs_away"`
	Flags     string            `json:"flags"`
	UniqueID  string            `json:"unique_id"`
}

type PriceUpdate struct {
	LpNum string `json:"lp_num"`
	LpDen string `json:"lp_den"`
}
