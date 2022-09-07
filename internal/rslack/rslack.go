package rslack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/rusq/slackinvite/internal/chttp"
)

const (
	domain  = "https://slack.com"
	baseURL = domain + "/api/"
)

type TeamInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	token string
	cl    *http.Client
}

// Response is the Slack server response.
type Response struct {
	Ok    bool     `json:"ok"`
	Error string   `json:"error"`
	Team  TeamInfo `json:"team,omitempty"`
}

func New(token string, cookies []*http.Cookie) *Client {
	cl := chttp.New(domain, cookies, nil)
	return &Client{token: token, cl: cl}
}

func NewDcookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:    "d",
		Value:   value,
		Path:    "/",
		Domain:  ".slack.com",
		Expires: time.Now().AddDate(10, 0, 0),
		Secure:  true,
	}
}

func (c *Client) AdminUsersInvite(teamID, email string) error {
	values := url.Values{
		"token":   {c.token},
		"email":   {email},
		"team_id": {teamID},
	}
	resp, err := c.cl.PostForm(baseURL+"admin.users.invite", values)
	if err != nil {
		return fmt.Errorf("admin.users.invite error: %s", err)
	}
	defer resp.Body.Close()
	if _, err := decode(resp); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetTeamInfo() (ti TeamInfo, err error) {
	values := url.Values{
		"token": {c.token},
	}
	resp, err := c.cl.Get(baseURL + "team.info?" + values.Encode())
	if err != nil {
		return ti, err
	}
	r, err := decode(resp)
	if err != nil {
		return ti, err
	}
	ti = r.Team
	return
}

func decode(resp *http.Response) (sr Response, err error) {
	if resp.StatusCode != http.StatusOK {
		return sr, fmt.Errorf("server returned NOT OK (code=%d slack error, if any: %s)", resp.StatusCode, sr.Error)
	}
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return sr, fmt.Errorf("error decoding slack response: %s", err)
	}
	if !sr.Ok || sr.Error != "" {
		return sr, fmt.Errorf("server returned an error (if any: %s)", sr.Error)
	}
	return
}
