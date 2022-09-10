package slackinvite

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
	"github.com/rusq/slackinvite/internal/chtml"
	"github.com/rusq/slackinvite/internal/recaptcha"
	"github.com/rusq/slackinvite/internal/rslack"
)

const (
	secretSz     = 64
	tokenTimeout = 20 * time.Minute // enough for time to think
)

var (
	//go:embed templates/*.html
	fs         embed.FS
	indexTmpl  = template.Must(chtml.NewLayout().WithLayout("templates/layout.html").ParseFS(fs, "templates/index.html"))
	thanksTmpl = template.Must(chtml.NewLayout().WithLayout("templates/layout.html").ParseFS(fs, "templates/thanks.html"))
)

// Server is the http server that issues the invites.
type Server struct {
	addr   string         // address to listen to
	teamID string         // slack team id
	secret [secretSz]byte // secret is the key used to encrypt the CSRF token.

	client slackClient // slack client that is called to issue the invite.
	db     *sql.DB     // (UNUSED) database connection
	rc     recaptcha.ReCaptcha

	fld Fields // template fields
}

type slackClient interface {
	GetTeamInfo() (rslack.TeamInfo, error)
	AdminUsersInvite(teamName string, emailAddress string) error
}

func New(addr string, db *sql.DB, client slackClient, rc recaptcha.ReCaptcha, fields Fields) (*Server, error) {
	var secret [secretSz]byte
	if _, err := io.ReadFull(rand.Reader, secret[:]); err != nil {
		return nil, err
	}
	ti, err := client.GetTeamInfo()
	if err != nil {
		return nil, err
	}

	s := &Server{
		client: client,
		db:     db,
		addr:   addr,
		secret: secret,
		teamID: ti.ID,
		rc:     rc,
		fld:    fields,
	}
	return s, nil
}

func (s *Server) Run() error {
	if err := initSecure("CSRF", 1024); err != nil {
		return err
	}

	dlog.Printf("Running inviter for %s, with team_id=%q", s.fld.SlackWorkspace, s.teamID)

	middleware.RequestIDHeader = "X-Request-ID" // https://github.com/heroku/x/blob/v0.0.52/requestid/requestid.go#L11

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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
