package main

import (
	"bytes"
	"context"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	firebase "firebase.google.com/go"
	authenticator "firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"strings"

	"encoding/json"
	"fmt"
	"github.com/rnpridgeon/cops/license"
	"github.com/rnpridgeon/exzenzd/zendesk"
)

var (
	app    *firebase.App
	client *authenticator.Client
	appCtx context.Context
	zd     *zendesk.Client
	clerk  *license.Agent
)

const Origin = "https://cops-onboarding.herokuapp.com"

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

func RootOpts(ctx *fasthttp.RequestCtx) {
	fmt.Println(ctx.Request.String())
	ctx.Response.Header.Add("Access-Control-Allow-Origin", Origin)
	ctx.Response.Header.Add("Access-Control-Allow-Headers", "*")
	ctx.Response.Header.Add("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
}

func Index(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")

	if Authenticated(ctx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody([]byte("Welcome!\n"))
	}

}

func CreateOrganization(ctx *fasthttp.RequestCtx) {
	if Authenticated(ctx) {
		//Needed to ensure object name is encoded in payload
		var record struct {
			Organization *zendesk.Organization `json:"organization"`
		}

		log.Printf("%s\n", ctx.Request.Body())

		if json.Unmarshal(ctx.Request.Body(), &record) == nil {
			log.Printf("INFO: Processing Create Organization request for %s\n", record.Organization.Name)
			// ZD organization names have to unique, this should suffice for the 'customer id' for now
			token, _ := clerk.SignToken(license.NewToken(record.Organization.Name,
				record.Organization.OrganizationFields.EffectiveAt.Unix(),
				record.Organization.OrganizationFields.RenewalAt.Unix(),
				true))

			record.Organization.OrganizationFields.LicenseKey = token
			enhanced, _ := json.Marshal(&record)

			zd.Create("organizations", enhanced, &ctx.Response)
			ctx.Response.Header.Set("Access-Control-Allow-Origin", Origin)

			return
		}
		// Json unmarshal failed, return bad request and malformed request body
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, "Failed to deserialize request %s\n", ctx.Request.Body())
	}
}

func CreateUser(ctx *fasthttp.RequestCtx) {
	if Authenticated(ctx) {
		log.Printf("%s\n", ctx.Request.Body())
		zd.Create("users", ctx.Request.Body(), &ctx.Response)
	}
}

func main() {
	var err error

	appCtx = context.Background()
	app, err = firebase.NewApp(appCtx, nil, option.WithCredentialsFile("exclude/firebase-conf.json"))
	client, err = app.Auth(appCtx)

	zd = zendesk.NewDefaultClient(zendesk.LoadCredentialsFile("exclude/zendesk-conf.json"))

	clerk, err = license.NewAgent("exclude/private_key.der", "exclude/public_key.der")

	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	router := fasthttprouter.New()
	router.GET("/", Index)
	router.POST("/organizations", CreateOrganization)
	router.POST("/users", CreateUser)
	router.OPTIONS("/*opts", RootOpts)
	log.Fatal(fasthttp.ListenAndServeTLS(":8080", "exclude/server.crt", "exclude/server.key", router.Handler))
	
	}
