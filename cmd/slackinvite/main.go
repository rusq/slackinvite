package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rusq/slackinviter"
	"github.com/rusq/slackinviter/internal/recaptcha"
	"github.com/slack-go/slack"
)

var _ = godotenv.Load()

const (
	addr = ":8080"
)

type params struct {
	Token  string
	Cookie string
	Addr   string
	RC     recaptcha.ReCaptcha
}

var cli params

func init() {
	flag.StringVar(&cli.Token, "t", os.Getenv("TOKEN"), "slack `token`")
	flag.StringVar(&cli.Cookie, "c", os.Getenv("COOKIE"), "slack `cookie`")
	flag.StringVar(&cli.Addr, "l", addr, "listener `address`")
	flag.StringVar(&cli.RC.SiteKey, "site-key", os.Getenv("RECAPTCHA_KEY"), "recaptcha `key`")
	flag.StringVar(&cli.RC.SecretKey, "site-secret", os.Getenv("RECAPTCHA_SECRET"), "recaptcha `secret`")
}

func main() {
	flag.Parse()
	if cli.Token == "" || cli.Cookie == "" {
		flag.Usage()
		log.Fatal("token or cookie not present")
	}
	if cli.Addr == "" {
		cli.Addr = addr
	}

	log.Printf("listening on %s", cli.Addr)

	client := slack.New(cli.Token, slack.OptionCookie("d", cli.Cookie))

	si, err := slackinviter.New(cli.Addr, nil, client, cli.RC, "Slackdump")
	if err != nil {
		log.Fatal(err)
	}
	if err := si.Run(); err != nil {
		log.Fatal(err)
	}
}
