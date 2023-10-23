package statistics

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const genderStatsURL = "https://api.genderize.io"

func (p *Provider) GenderByName(name string) (gender string, err error) {
	url := fmt.Sprintf("%s/?name=%s", genderStatsURL, name)
	var r *http.Response
	r, err = http.Get(url)

	if err != nil {
		return "", fmt.Errorf("failed to receive gender stats (%s): %w", url, err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to receive gender stats (%s): status %d", url, r.StatusCode)
	}

	data := &struct{ Gender string }{}

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return "", fmt.Errorf("failed to deserialize gender stats response (%s): %w", url, err)
	}

	return data.Gender, nil
}
