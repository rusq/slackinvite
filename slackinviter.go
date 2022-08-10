package slackinviter

import (
	"crypto/rand"
	"database/sql"
	"embed"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rusq/secure"
	"github.com/slack-go/slack"
)

const (
	secretSz     = 64
	tokenTimeout = 20 * time.Minute // enough for time to think
)

//go:embed templates/*.html
var fs embed.FS

var tmpl = template.Must(template.ParseFS(fs, "templates/*.html"))

type Server struct {
	name   string // title to show on the page
	addr   string
	teamID string
	secret [secretSz]byte
	client *slack.Client
	db     *sql.DB
}

func New(addr string, db *sql.DB, client *slack.Client, title string) (*Server, error) {
	var secret [secretSz]byte
	if _, err := io.ReadFull(rand.Reader, secret[:]); err != nil {
		return nil, err
	}
	ti, err := client.GetTeamInfo()
	if err != nil {
		return nil, err
	}

	s := &Server{name: title, client: client, db: db, addr: addr, secret: secret, teamID: ti.ID}
	return s, nil
}

func (s *Server) Run() error {
	secure.SetSignature("CSRF")

	log.Printf("Running inviter for %s, with team_id=%q", s.name, s.teamID)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/", http.HandlerFunc(s.handleRoot))
	r.Handle("/thankyou", http.HandlerFunc(s.handleThankyou))

	return http.ListenAndServe(s.addr, r)
}

type page struct {
	Title         string
	Name          string
	Token         string
	Email         string
	SubmitBtnText string
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
