package slackinviter

import (
	"crypto/rand"
	"database/sql"
	"embed"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rusq/dlog"
	"github.com/rusq/secure"
	"github.com/rusq/slackinviter/internal/recaptcha"
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
	rc     recaptcha.ReCaptcha
}

func New(addr string, db *sql.DB, client *slack.Client, rc recaptcha.ReCaptcha, title string) (*Server, error) {
	var secret [secretSz]byte
	if _, err := io.ReadFull(rand.Reader, secret[:]); err != nil {
		return nil, err
	}
	ti, err := client.GetTeamInfo()
	if err != nil {
		return nil, err
	}

	s := &Server{
		name:   title,
		client: client,
		db:     db,
		addr:   addr,
		secret: secret,
		teamID: ti.ID,
		rc:     rc,
	}
	return s, nil
}

func (s *Server) Run() error {
	if err := initSecure("CSRF", 1024); err != nil {
		return err
	}

	dlog.Printf("Running inviter for %s, with team_id=%q", s.name, s.teamID)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/", http.HandlerFunc(s.handleRoot))
	r.Handle("/thankyou", http.HandlerFunc(s.handleThankyou))

	return http.ListenAndServe(s.addr, r)
}

func initSecure(sig string, saltSz int) error {
	secure.SetSignature(sig)

	var salt = make([]byte, saltSz)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}
	secure.SetSalt(salt)

	return nil
}
