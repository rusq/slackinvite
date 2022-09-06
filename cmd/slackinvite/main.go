package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rusq/dlog"
	"github.com/rusq/osenv/v2"
	si "github.com/rusq/slackinvite"
	"github.com/slack-go/slack"

	"github.com/rusq/slackinvite/internal/recaptcha"
)

var _ = godotenv.Load()

const (
	defPort = "8080"
)

type params struct {
	Token     string
	Cookie    string
	Addr      string
	Port      string
	RC        recaptcha.ReCaptcha
	FieldsCfg string
}

var cmdline params

func init() {
	flag.StringVar(&cmdline.FieldsCfg, "cfg", os.Getenv("CONFIG_FILE"), "Config file with template values")
	flag.StringVar(&cmdline.Token, "t", os.Getenv("TOKEN"), "slack `token`")
	flag.StringVar(&cmdline.Cookie, "c", os.Getenv("COOKIE"), "slack `cookie`")
	flag.StringVar(&cmdline.Addr, "addr", os.Getenv("ADDR"), "host `address` for the listener, i.e. 127.0.0.1, if empty\nwill listen on all interfaces.")
	flag.StringVar(&cmdline.Port, "port", osenv.Value("PORT", defPort), "`port` to listen to")
	flag.StringVar(&cmdline.RC.SiteKey, "site-key", os.Getenv("RECAPTCHA_KEY"), "recaptcha `key`")
	flag.StringVar(&cmdline.RC.SecretKey, "site-secret", os.Getenv("RECAPTCHA_SECRET"), "recaptcha `secret`")
}

func main() {
	flag.Parse()

	if cmdline.Token == "" || cmdline.Cookie == "" {
		flag.Usage()
		dlog.Fatal("token or cookie not present")
	}
	if cmdline.Port == "" {
		cmdline.Port = defPort
	}

	client := slack.New(cmdline.Token, slack.OptionCookie("d", cmdline.Cookie))

	fields, err := si.LoadFields(cmdline.FieldsCfg)
	if err != nil {
		log.Fatal(err)
	}

	listenerAddr := cmdline.Addr + ":" + cmdline.Port
	dlog.Printf("listening on %s", listenerAddr)
	si, err := si.New(listenerAddr, nil, client, cmdline.RC, fields)
	if err != nil {
		dlog.Fatal(err)
	}
	if err := si.Run(); err != nil {
		dlog.Fatal(err)
	}
}
