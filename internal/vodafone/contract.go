package vodafone

type GuestResponse struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Result      struct {
		Token                  string `json:"token"`
		WidgetAutoRefreshDelay string `json:"widgetAutoRefreshDelay"`
		WidgetNextRefreshDelay string `json:"widgetNextRefreshDelay"`
	} `json:"result"`
}

type LoginResponse struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Result      struct {
		Customer struct {
			Email      string `json:"email"`
			Firstname  string `json:"firstname"`
			FiscalCode string `json:"fiscalCode"`
			Lastname   string `json:"lastname"`
			TypeID     string `json:"typeId"`
			Username   string `json:"username"`
		} `json:"customer"`
		Items []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"items"`
		PartyIdentifier string `json:"partyIdentifier"`
		PsgToken        string `json:"psgToken"`
	} `json:"result"`
}

type Counter struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Threshold struct {
		Values []struct {
			Value             float64 `json:"value"`
			Threshold         float64 `json:"threshold"`
			Icon              string  `json:"icon"`
			ThresholdMessage  string  `json:"thresholdMessage"`
			Order             int     `json:"order"`
			Unlimited         bool    `json:"unlimited"`
			Unit              string  `json:"unit"`
			UnitExtended      string  `json:"unitExtended"`
			DirectionCategory string  `json:"directionCategory"`
			ResidualValue     float64 `json:"residual_value"`
			PeriodStart       string  `json:"period_start"`
			PeriodEnd         string  `json:"period_end"`
			GhieraLabel       string  `json:"ghiera_label"`
			WidgetLabel       string  `json:"widget_label"`
			ShowThreshold     bool    `json:"show_threshold"`
			IsPassion         bool    `json:"isPassion"`
			IsResidual        bool    `json:"is_residual"`
		} `json:"values"`
		Size          int  `json:"size"`
		ShowThreshold bool `json:"show_threshold"`
	} `json:"threshold"`
	HomeVisible    bool   `json:"homeVisible"`
	CounterMessage string `json:"counterMessage"`
	Order          int    `json:"order"`
	Aggregated     bool   `json:"aggregated"`
	ID             string `json:"id"`
	PeriodStart    string `json:"period_start"`
	PeriodEnd      string `json:"period_end"`
	Banner         struct {
		Image      string `json:"image"`
		Text       string `json:"text"`
		AnchorURL  string `json:"anchorUrl"`
		AnchorText string `json:"anchorText"`
	} `json:"banner,omitempty"`
	DynamicResetLabel string `json:"dynamic_reset_label,omitempty"`
}
type CountersResponse struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Result      struct {
		Counters []Counter `json:"counters"`
	} `json:"result"`
}
