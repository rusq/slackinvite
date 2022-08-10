module github.com/rusq/slackinviter

go 1.18

require (
	github.com/go-chi/chi/v5 v5.0.7
	github.com/joho/godotenv v1.4.0
	github.com/rusq/secure v0.0.4
	github.com/slack-go/slack v0.11.2
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
)

replace github.com/slack-go/slack => github.com/rusq/slack v0.11.100
