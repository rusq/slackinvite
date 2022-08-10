package slackinviter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type gResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func verifyRecaptcha(secret string, response string) error {
	values := url.Values{
		"secret":   []string{secret},
		"response": []string{response},
	}

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var gr gResponse
	if err := dec.Decode(&gr); err != nil {
		return err
	}
	if !gr.Success {
		return fmt.Errorf("recapcha at %s on %q is unsuccessful: %v", gr.ChallengeTs, gr.Hostname, gr.ErrorCodes)
	}
	return nil
}
