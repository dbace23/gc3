package rapid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BMIResponse struct {
	Bmi          float64 `json:"bmi"`
	Health       string  `json:"health"`
	HealthyRange string  `json:"healthy_bmi_range"`
}

type BMIClient struct {
	Key  string
	Host string
	HC   *http.Client
}

func NewBMI(key, host string) *BMIClient {
	return &BMIClient{
		Key:  key,
		Host: host,
		HC:   &http.Client{Timeout: 6 * time.Second},
	}
}

func (c *BMIClient) Calculate(weightKg, heightCm int) (*BMIResponse, error) {
	payload := map[string]int{"weight": weightKg, "height": heightCm}
	b, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", fmt.Sprintf("https://%s/metric", c.Host), bytes.NewReader(b))
	req.Header.Set("x-rapidapi-key", c.Key)
	req.Header.Set("x-rapidapi-host", c.Host)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HC.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("bmi api error: %s", res.Status)
	}
	var out BMIResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
