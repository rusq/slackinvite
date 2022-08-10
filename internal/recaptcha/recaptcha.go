package recaptcha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const htmlID = "html_element"

type ReCaptcha struct {
	SiteKey   string
	SecretKey string
}

type gResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func (rc ReCaptcha) Verify(response string) error {
	values := url.Values{
		"secret":   []string{rc.SecretKey},
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

func (rc ReCaptcha) JS() string {
	return `<script type="text/javascript">
	var onloadCallback = function() {
		grecaptcha.render('` + htmlID + `', {
			'sitekey' : '` + rc.SiteKey + `'
			});
			};
			</script>`
}

func (rc ReCaptcha) HTML() string {
	return `<div id="` + htmlID + `"></div>`
}
