module github.com/rusq/slackinvite

go 1.23.0

toolchain go1.24.2

require (
	github.com/go-chi/chi/v5 v5.2.1
	github.com/goccy/go-yaml v1.18.0
	github.com/joho/godotenv v1.5.1
	github.com/rusq/chttp v1.1.0
	github.com/rusq/dlog v1.4.0
	github.com/rusq/osenv/v2 v2.0.1
	github.com/rusq/secure v0.0.4
)

require (
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/net v0.41.0 // indirect
)

// replace github.com/slack-go/slack => github.com/rusq/slack v0.11.100
