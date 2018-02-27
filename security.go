package main

import (
	"bytes"
	"context"
	"log"
	"strings"

	firebase "firebase.google.com/go"
	authenticator "firebase.google.com/go/auth"
	"google.golang.org/api/option"

	"github.com/valyala/fasthttp"
)

var (
	app    *firebase.App
	client *authenticator.Client
)

func getGroup(subject *authenticator.UserRecord) (group string) {
	if subject != nil {
		return strings.SplitN(subject.UserInfo.Email, "@", -1)[1]
	}

	log.Printf("ERROR: Unable to process subject %v\n", subject)

	return ""
}

// TODO: set up custom claims to make role auth less hacky
func Authenticated(ctx *fasthttp.RequestCtx) (authenticated bool) {
	var subject *authenticator.UserRecord
	authenticated = false

	// Authorization: bearer authentication_token
	authz := bytes.Split(ctx.Request.Header.Peek("Authorization"), []byte("bearer"))

	// Invalid authorization header: deny access
	if len(authz) != 2 {

		ctx.SetStatusCode(fasthttp.StatusForbidden)
		ctx.SetBody([]byte("Access Denied"))

		log.Printf("ERROR: requester at %v used invalid authorization header; denying request",
			ctx.RemoteAddr())

		return authenticated
	}

	// VerifyIDToken returns an error for invalid tokens, if its true check group
	user, err := client.VerifyIDToken(string(bytes.TrimSpace(authz[1])))

	if err != nil {
		authenticated = false
		return authenticated
	}

	subject, _ = client.GetUser(appCtx, user.UID)
	authenticated = getGroup(subject) == "confluent.io"
	log.Printf("INFO: %v@%v requested %s on %s; result: %t\n", subject.Email, ctx.RemoteAddr(), ctx.Method(),
		ctx.URI().Path(), authenticated)

	return authenticated
}

func InitAuthenticator(confFile string) {
	appCtx = context.Background()

	app, err := firebase.NewApp(appCtx, nil, option.WithCredentialsFile(confFile))
	client, err = app.Auth(appCtx)

	if err != nil {
		log.Fatal("Failed to configure Authenticator, shtting down")
	}
}
