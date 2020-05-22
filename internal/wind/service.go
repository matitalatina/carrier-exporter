package wind

import "net/http"

type Service struct {
	Client      Client
	authCookies []*http.Cookie
}

func (s *Service) GetInsights(credentials Credentials, lineId string, contractId string) (*InsightsSummary, error) {
	authCookies, err := s.fetchAuthCookies(credentials)

	if err != nil {
		return nil, err
	}

	stats, err := s.Client.GetStats(authCookies, lineId, contractId)

	if err != nil {
		return nil, err
	}

	// TODO: it takes the first one. We can improve
	return &stats.Lines[0].InsightsSummary, nil
}

func (s *Service) fetchAuthCookies(credentials Credentials) ([]*http.Cookie, error) {
	if s.authCookies != nil {
		return s.authCookies, nil
	}

	authCookies, err := s.Client.Login(credentials)

	if err != nil {
		return nil, err
	}

	s.authCookies = authCookies

	return authCookies, nil
}
