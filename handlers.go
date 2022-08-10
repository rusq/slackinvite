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
	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching
	w.Header().Add("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
	pg := page{
		Title:       "Join " + s.fld.SlackWorkspace + " Slack",
		Token:       csrf,
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
	email := r.FormValue("email")
	csrf := r.FormValue("token")
	recaptcha := r.FormValue("g-recaptcha-response")
	dlog.Debugf("%#v", r.Form)
	if csrf == "" {
		dlog.Print("empty token")
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	if err := verifyToken(csrf, s.secret, tokenTimeout); err != nil {
		dlog.Print(err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		dlog.Print(err)
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	resp, err := s.rc.Verify(recaptcha)
	if err != nil {
		dlog.Println(err)
		http.Error(w, "invalid captcha", http.StatusExpectationFailed)
		return
	}
	if resp.Score != "" {
		score, err := resp.Score.Float64()
		if err != nil {
			dlog.Printf("failed to interpret score: %s", resp.Score)
		} else if score < 0.5 {
			dlog.Printf("low re-captcha score for %s: %s", email, resp.Score)
			http.Error(w, "403 you are known for suspicious behaviour, please leave", http.StatusForbidden)
			return
		}
	}

	// if err := s.client.InviteToTeam(s.teamID, "Test", "Invite", email); err != nil {
	// 	dlog.Printf("email: %s: %s", email, err)
	// 	http.Error(w, "something went wrong", http.StatusInternalServerError)
	// 	return
	// }
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
