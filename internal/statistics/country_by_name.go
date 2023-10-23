package statistics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const countryStatsURL = "https://api.nationalize.io"

type StatsDataCountry struct {
	Country []*StatsDataCountryId
}

type StatsDataCountryId struct {
	CountryId string `json:"country_id"`
}

func (p *Provider) CountryByName(name string) (country string, err error) {
	url := fmt.Sprintf("%s/?name=%s", countryStatsURL, url.QueryEscape(name))
	var r *http.Response
	r, err = http.Get(url)

	if err != nil {
		return "", fmt.Errorf("failed to receive country stats (%s): %w", url, err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to receive country stats (%s): status %d", url, r.StatusCode)
	}

	data := &StatsDataCountry{Country: make([]*StatsDataCountryId, 0)}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return "", fmt.Errorf("failed to deserialize country stats response (%s): %w", url, err)
	}

	if len(data.Country) > 0 {
		return data.Country[0].CountryId, nil
	}

	return "", nil
}
