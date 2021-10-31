package wind

type Service struct {
	Client    Client
	authToken string
}

func (s *Service) GetInsights(credentials Credentials, lineId string, contractId string) (*InsightsSummary, error) {
	authToken, err := s.fetchAuthToken(credentials)

	if err != nil {
		return nil, err
	}

	stats, err := s.Client.GetStats(authToken, lineId, contractId)

	if err != nil {
		return nil, err
	}

	// TODO: it takes the first one. We can improve
	return &stats.Data.Lines[0].InsightsSummary, nil
}

func (s *Service) fetchAuthToken(credentials Credentials) (string, error) {
	if s.authToken != "" {
		return s.authToken, nil
	}

	authCookies, err := s.Client.Login(credentials)

	if err != nil {
		return "", err
	}

	s.authToken = authCookies

	return authCookies, nil
}
