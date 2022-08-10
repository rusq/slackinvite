package recaptcha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/rusq/dlog"
)

type ReCaptcha struct {
	SiteKey   string
	SecretKey string
}

type Response struct {
	Success     bool        `json:"success,omitempty"`
	Score       json.Number `json:"score,omitempty"`
	Action      string      `json:"action,omitempty"`
	ChallengeTs string      `json:"challenge_ts,omitempty"`
	Hostname    string      `json:"hostname,omitempty"`
	ErrorCodes  []string    `json:"error-codes,omitempty"`
}

const validationJS = `<script>(function () {
	'use strict'
  
	window.addEventListener('load', function () {
	  // Fetch all the forms we want to apply custom Bootstrap validation styles to
	  var forms = document.getElementsByClassName('needs-validation')
  
	  // Loop over them and prevent submission
	  Array.prototype.filter.call(forms, function (form) {
		form.addEventListener('submit', function (event) {
		  if (form.checkValidity() === false) {
			event.preventDefault()
			event.stopPropagation()
		  }
		  form.classList.add('was-validated')
		}, false)
	  })
	}, false)
  }())</script>`

func (rc ReCaptcha) Verify(response string) (*Response, error) {
	values := url.Values{
		"secret":   []string{rc.SecretKey},
		"response": []string{response},
	}

	apiResp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", values)
	if err != nil {
		return nil, err
	}
	defer apiResp.Body.Close()
	dec := json.NewDecoder(apiResp.Body)
	var resp Response
	if err := dec.Decode(&resp); err != nil {
		return nil, err
	}
	dlog.Printf("%#v", resp)
	if !resp.Success {
		return &resp, fmt.Errorf("recapcha at %s on %q is unsuccessful: %v", resp.ChallengeTs, resp.Hostname, resp.ErrorCodes)
	}
	return &resp, nil
}

// JSv2 returns the javascript code that should be placed in the HTML template.
// Make sure that htmlID matches the htmlID passed to HTML().
func (rc ReCaptcha) JSv2(htmlID string) string {
	return `<script type="text/javascript">
	var onloadCallback = function() {
		grecaptcha.render('` + htmlID + `', {
			'sitekey' : '` + rc.SiteKey + `'
			});
		};
		</script>
		<script src="https://www.google.com/recaptcha/api.js?onload=onloadCallback&render=explicit" async defer></script>
	`
}

// HTMLv2 returns the HTMLv2 code that should be placed inside the <form>, possibly
// next to the "submit" button.  Make sure that htmlID matches the htmlID passed
// to JS.
func (rc ReCaptcha) HTMLv2(htmlID string) string {
	return `<div id="` + htmlID + `"></div>`
}

// JSv3 returns the javascript code necessary for reCaptcha v3.  formID should
// match the <form id="<value>">.
func (rc ReCaptcha) JSv3(formID string) string {
	return validationJS + `<script src="https://www.google.com/recaptcha/api.js"></script>
	<script>
	function onSubmit(token) {
	  document.getElementById("` + formID + `").submit();
	}
  </script>`

}

// HTMLv3 returns the HTML code necessary for reCaptcha v3.  It should be
// placed instead of the "Submit" button, as it returns the code for it.
// buttonText is the text of the button, add any additional classes you want
// the button to have by providing classes.
// Drawbacks:  form validation should be handled manually.
func (rc ReCaptcha) HTMLv3(buttonText string, classes ...string) string {
	return `<button class="g-recaptcha ` + strings.Join(classes, " ") + `" 
	data-sitekey="` + rc.SiteKey + `" 
	data-callback='onSubmit' 
	data-action='submit'>` + buttonText + `</button>`
}
