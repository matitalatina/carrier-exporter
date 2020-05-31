package vodafone

import (
	"net/http"
)

type Service struct {
	Client      Client
	authCookies []*http.Cookie
}

func (s *Service) GetCounters(credentials Credentials, phone string) (*CountersResponse, error) {
	guestCredentials, err := s.Client.GetToken()
	apiAuth := ApiAuth{
		BwbSessionId: guestCredentials.BwbSessionId,
		Token:        guestCredentials.Response.Result.Token,
	}

	_, err = s.Client.Login(apiAuth, credentials)

	if err != nil {
		return nil, err
	}

	return s.Client.GetCounters(apiAuth, phone)
}
