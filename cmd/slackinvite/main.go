package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rusq/dlog"
	si "github.com/rusq/slackinviter"
	"github.com/rusq/slackinviter/internal/recaptcha"
	"github.com/slack-go/slack"
)

var _ = godotenv.Load()

const (
	addr = ":8080"
)

type params struct {
	Token     string
	Cookie    string
	Addr      string
	RC        recaptcha.ReCaptcha
	FieldsCfg string
}

var cmdline params

func init() {
	flag.StringVar(&cmdline.FieldsCfg, "cfg", os.Getenv("CONFIG_FILE"), "Config file with template values")
	flag.StringVar(&cmdline.Token, "t", os.Getenv("TOKEN"), "slack `token`")
	flag.StringVar(&cmdline.Cookie, "c", os.Getenv("COOKIE"), "slack `cookie`")
	flag.StringVar(&cmdline.Addr, "l", addr, "listener `address`")
	flag.StringVar(&cmdline.RC.SiteKey, "site-key", os.Getenv("RECAPTCHA_KEY"), "recaptcha `key`")
	flag.StringVar(&cmdline.RC.SecretKey, "site-secret", os.Getenv("RECAPTCHA_SECRET"), "recaptcha `secret`")
}

func main() {
	flag.Parse()
	if cmdline.Token == "" || cmdline.Cookie == "" {
		flag.Usage()
		dlog.Fatal("token or cookie not present")
	}
	if cmdline.Addr == "" {
		cmdline.Addr = addr
	}

	dlog.Printf("listening on %s", cmdline.Addr)

	client := slack.New(cmdline.Token, slack.OptionCookie("d", cmdline.Cookie))

	fields, err := si.LoadFields(cmdline.FieldsCfg)
	if err != nil {
		log.Fatal(err)
	}

	si, err := si.New(cmdline.Addr, nil, client, cmdline.RC, fields)
	if err != nil {
		dlog.Fatal(err)
	}
	if err := si.Run(); err != nil {
		dlog.Fatal(err)
	}
}
