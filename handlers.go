package slackinviter

import (
	"html/template"
	"net/http"
	"net/mail"

	"github.com/rusq/dlog"
	"github.com/rusq/secure"
)

type page struct {
	Title       string
	Token       string
	Email       string
	Error       string
	FormID      string
	CaptchaJS   template.HTML
	CaptchaHTML template.HTML
	Fields      Fields
}

type errorCode string

const (
	errCodeToken   errorCode = "SI001"
	errCodeEmail   errorCode = "SI002"
	errCodeCaptcha errorCode = "SI003"
	errCodeBot     errorCode = "SI004"
	errCodeSlack   errorCode = "SI005"
	errCodeServer  errorCode = "SI500"
)

var errorText = map[errorCode]string{
	errCodeToken:   "Invalid token.  Please refresh the page and try again.",
	errCodeEmail:   "Invalid email.  Make sure you're entering a correct email.",
	errCodeCaptcha: "Invalid captcha. Make sure you're not a robot.",
	errCodeBot:     "You are known for being a robot.  Please leave.",
	errCodeSlack:   "Failed to send the invitation.",
	errCodeServer:  "Server error occurred.",
}

func (ec errorCode) String() string {
	if ec == "" {
		return ""
	}
	txt, ok := errorText[ec]
	if !ok {
		return "Unknown error"
	}
	return txt
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.hgetRoot(w, r)
	case http.MethodPost:
		s.hpostRoot(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *Server) hgetRoot(w http.ResponseWriter, r *http.Request) {
	const formID = "invite-form"
	csrf, err := generateToken(s.secret)
	if err != nil {
		dlog.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	pageErr := r.FormValue("e")

	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching
	w.Header().Add("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
	pg := page{
		Title:       "Join " + s.fld.SlackWorkspace + " Slack",
		Token:       csrf,
		Error:       errorCode(pageErr).String(),
		FormID:      formID,
		CaptchaJS:   template.HTML(s.rc.JSv3(formID)),
		CaptchaHTML: template.HTML(s.rc.HTMLv3(s.fld.SubmitButton, "btn", "btn-primary")),
		Fields:      s.fld,
	}
	if err := tmpl.ExecuteTemplate(w, "index.html", pg); err != nil {
		dlog.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) hpostRoot(w http.ResponseWriter, r *http.Request) {
	errCode := r.FormValue("e")
	if errCode != "" {
		s.hgetRoot(w, r)
		return
	}

	email := r.FormValue("email")
	csrf := r.FormValue("token")
	recaptcha := r.FormValue("g-recaptcha-response")
	dlog.Debugf("%#v", r.Form)
	if csrf == "" {
		dlog.Print("empty token")
		errRedirect(w, r, errCodeToken)
		return
	}

	if err := verifyToken(csrf, s.secret, tokenTimeout); err != nil {
		dlog.Print(err)
		errRedirect(w, r, errCodeToken)
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		dlog.Print(err)
		errRedirect(w, r, errCodeEmail)
		return
	}

	resp, err := s.rc.Verify(recaptcha)
	if err != nil {
		dlog.Println(err)
		errRedirect(w, r, errCodeCaptcha)
		return
	}
	if resp.Score != "" {
		score, err := resp.Score.Float64()
		if err != nil {
			dlog.Printf("failed to interpret score: %s", resp.Score)
		} else if score < 0.5 {
			dlog.Printf("low re-captcha score for %s: %s", email, resp.Score)
			errRedirect(w, r, errCodeBot)
			return
		}
	}

	if err := s.client.InviteToTeam(s.teamID, "Test", "Invite", email); err != nil {
		dlog.Printf("email: %s: %s", email, err)
		errRedirect(w, r, errCodeSlack)
		return
	}
	dlog.Printf("successfully invited: %q", email)
	ct, err := secure.EncryptWithPassphrase(email, s.secret[:])
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, "/thankyou?to="+ct, http.StatusMovedPermanently)
}

func (s *Server) handleThankyou(w http.ResponseWriter, r *http.Request) {
	ct := r.FormValue("to")
	if ct == "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	email, err := secure.DecryptWithPassphrase(ct, s.secret[:])
	if err != nil {
		dlog.Print(err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	pg := page{
		Title:  "Great success!",
		Email:  email,
		Fields: s.fld,
	}
	if err := tmpl.ExecuteTemplate(w, "thanks.html", pg); err != nil {
		dlog.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// errRedirect redirects to the root page and sets the error code to code.
func errRedirect(w http.ResponseWriter, r *http.Request, code errorCode) {
	http.Redirect(w, r, "/?e="+string(code), http.StatusTemporaryRedirect)
}
