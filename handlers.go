package slackinviter

import (
	"html/template"
	"log"
	"net/http"

	"github.com/rusq/secure"
)

type page struct {
	Title         string
	Name          string
	Token         string
	Email         string
	SubmitBtnText string
	Error         string
	CaptchaJS     template.HTML
	CaptchaHTML   template.HTML
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
	csrf, err := generateToken(s.secret)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching
	w.Header().Add("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
	pg := page{
		Title:         "Join " + s.name + " Slack",
		Name:          s.name,
		Token:         csrf,
		SubmitBtnText: "Gimme, gimme",
		CaptchaJS:     template.HTML(s.rc.JS()),
		CaptchaHTML:   template.HTML(s.rc.HTML()),
	}
	if err := tmpl.ExecuteTemplate(w, "index.html", pg); err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) hpostRoot(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	csrf := r.FormValue("token")
	recaptcha := r.FormValue("g-recaptcha-response")
	if csrf == "" {
		log.Print("empty token")
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	if err := verifyToken(csrf, s.secret, tokenTimeout); err != nil {
		log.Print(err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	if email == "" {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	if err := s.rc.Verify(recaptcha); err != nil {
		log.Println(err)
		http.Error(w, "invalid captcha", http.StatusExpectationFailed)
		return
	}

	// if err := s.client.InviteToTeam(s.teamID, "Test", "Invite", email); err != nil {
	// 	log.Printf("email: %s: %s", email, err)
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
		log.Print(err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	pg := page{
		Title: "Great success!",
		Name:  s.name,
		Email: email,
	}
	if err := tmpl.ExecuteTemplate(w, "thanks.html", pg); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
