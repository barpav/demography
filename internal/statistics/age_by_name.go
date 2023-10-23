package statistics

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const ageStatsURL = "https://api.agify.io"

func (p *Provider) AgeByName(name string) (age int, err error) {
	url := fmt.Sprintf("%s/?name=%s", ageStatsURL, name)
	var r *http.Response
	r, err = http.Get(url)

	if err != nil {
		return 0, fmt.Errorf("failed to receive age stats (%s): %w", url, err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to receive age stats (%s): status %d", url, r.StatusCode)
	}

	data := &struct{ Age int }{}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return 0, fmt.Errorf("failed to deserialize age stats response (%s): %w", url, err)
	}

	return data.Age, nil
}
