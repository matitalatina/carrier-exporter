package tim

import (
	"errors"
	"strconv"
)

type Service struct {
	Client Client
}

func (s *Service) GetAggregate(credentials Credentials, phone string) (*AggregateResponse, error) {
	resp, err := s.Client.Authorize()

	if err != nil {
		return nil, err
	}

	resp, err = s.Client.Login(*resp, credentials)

	if err != nil {
		return nil, err
	}

	code, err := s.Client.Consent(*resp)

	if err != nil {
		return nil, err
	}

	authToken, err := s.Client.Token(*code)

	if err != nil {
		return nil, err
	}

	aggregate, err := s.Client.Aggregate(AuthApi{
		SessionID:   resp.SessionID,
		AccessToken: authToken.AccessToken,
	}, phone)

	if err != nil {
		return nil, err
	}

	return aggregate, nil
}

func (s *Service) GetAvailableBytes(credentials Credentials, phone string) (*float64, error) {
	aggregate, err := s.GetAggregate(credentials, phone)

	if err != nil {
		return nil, err
	}

	for _, a := range aggregate.Aggregateoffers {
		if a.Type == "DATI" {
			availableGb, err := strconv.ParseFloat(a.Value, 64)
			if err != nil {
				return nil, err
			}

			availableBytes := availableGb * 1024 * 1024 * 1024
			return &availableBytes, nil
		}
	}

	return nil, errors.New("Data not found")
}
