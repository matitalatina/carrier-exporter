package wind

type InsightsSummary struct {
	National struct {
		Data struct {
			Total     float64 `json:"total"`
			Available float64 `json:"available"`
			Unlimited bool    `json:"unlimited"`
		} `json:"data"`
	} `json:"national"`
	AllOfferSuspended bool `json:"allOfferSuspended"`
}

type Stats struct {
	Data struct {
		ID          string `json:"id"`
		CustomerID  string `json:"customerId"`
		Status      string `json:"status"`
		PaymentType string `json:"paymentType"`
		Sme         bool   `json:"sme"`
		Lines       []struct {
			Alias                      string   `json:"alias"`
			ID                         string   `json:"id"`
			LineAdslID                 string   `json:"lineAdslId"`
			CustomerID                 string   `json:"customerId"`
			ContractID                 string   `json:"contractId"`
			ActivationDate             string   `json:"activationDate"`
			Mobile                     bool     `json:"mobile"`
			Sme                        bool     `json:"sme"`
			PaymentType                string   `json:"paymentType"`
			Features                   []string `json:"features"`
			FeaturesDescription        string   `json:"featuresDescription"`
			LockedCredentialsRetrieval bool     `json:"lockedCredentialsRetrieval"`
			LockedAutoLoginFromApps    bool     `json:"lockedAutoLoginFromApps"`
			SimCard                    struct {
				LineID  string `json:"lineId"`
				SimType string `json:"simType"`
				Puk     string `json:"puk"`
			} `json:"simCard"`
			TariffPlan struct {
				Code           string        `json:"code"`
				Name           string        `json:"name"`
				Status         string        `json:"status"`
				ActivationDate string        `json:"activationDate"`
				RenewalDate    string        `json:"renewalDate"`
				ActionLinks    []interface{} `json:"actionLinks"`
				Channels       []string      `json:"channels"`
				Master         string        `json:"master"`
				Peso           string        `json:"peso"`
				TokenValue     []string      `json:"tokenValue"`
				OfferDetailTre []interface{} `json:"offerDetailTre"`
			} `json:"tariffPlan"`
			Options []struct {
				Code           string        `json:"code"`
				Name           string        `json:"name"`
				Status         string        `json:"status"`
				ActivationDate string        `json:"activationDate"`
				RenewalDate    string        `json:"renewalDate"`
				ActionLinks    []interface{} `json:"actionLinks"`
				Channels       []string      `json:"channels"`
				Master         string        `json:"master"`
				Peso           string        `json:"peso"`
				TokenValue     []string      `json:"tokenValue"`
				Insights       []struct {
					Type          string  `json:"type"`
					UnitOfMeasure string  `json:"unitOfMeasure"`
					Available     float64 `json:"available"`
					Total         float64 `json:"total"`
					Unlimited     bool    `json:"unlimited"`
					Group         string  `json:"group"`
				} `json:"insights"`
				Family struct {
					RootIntegrationID string `json:"rootIntegrationId"`
					OfferCode         string `json:"offerCode"`
					FamilyID          string `json:"familyId"`
					FamilyCode        string `json:"familyCode"`
				} `json:"family"`
				TreOptionCanDeactivate bool `json:"treOptionCanDeactivate"`
			} `json:"options"`
			Promos           []interface{}   `json:"promos"`
			InsightsSummary  InsightsSummary `json:"insightsSummary"`
			IntegrationStack string          `json:"integrationStack"`
			GracePeriod      struct {
			} `json:"gracePeriod"`
			DataAttivazionePhoenix string `json:"dataAttivazionePhoenix"`
			LineStatusInfo         string `json:"lineStatusInfo"`
			NextEventDate          struct {
				Date          string `json:"date"`
				DataEventType string `json:"dataEventType"`
				RopzCode      string `json:"ropzCode"`
				Tied          bool   `json:"tied"`
			} `json:"nextEventDate"`
		} `json:"lines"`
		LineType       string `json:"lineType"`
		ShowMGM        string `json:"showMGM"`
		ScontoFee      bool   `json:"scontoFee"`
		FlagMigrazione string `json:"flagMigrazione"`
	} `json:"data"`
	Status     string        `json:"status"`
	ErrorCodes []interface{} `json:"errorCodes"`
	Messages   []interface{} `json:"messages"`
}
