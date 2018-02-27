package main

import (
	"context"
	"log"
	"strconv"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/rnpridgeon/exzenzd/datastore"
	"github.com/rnpridgeon/utils/configuration"
)

var (
	appCtx context.Context
	ds     *datastore.Datasource
)

type Config struct {
	//ZDconf *zendesk.ZendeskConfig `json:"zendesk"`
	DBconf *datastore.Config `json:"database"`
}

const Origin = " https://cops-tooling.herokuapp.com"

func ConfigCORS(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("Access-Control-Allow-Origin", Origin)
	ctx.Response.Header.Add("Access-Control-Allow-Headers", "*")
	ctx.Response.Header.Add("Access-Control-Allow-Methods", "OPTIONS, GET")
}

func Index(ctx *fasthttp.RequestCtx) {
	if Authenticated(ctx) {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody([]byte("Welcome!\n"))
	}
}

func GetOrganizations(ctx *fasthttp.RequestCtx) {
	if Authenticated(ctx) {
		ds.ExportOrganizations(ctx)
	}
}

func GetOrganizationTickets(ctx *fasthttp.RequestCtx) {
	if Authenticated(ctx) {
		orgid, _ := strconv.ParseInt(ctx.UserValue("id").(string), 10, 64)
		ds.ExportOrganizationTickets(orgid, ctx)
	}
}

func main() {
	var err error
	var conf Config

	//configuration.FromFile(os.Args[1])(&conf)
	configuration.FromFile("./exclude/conf.json")(&conf)

	InitAuthenticator("exclude/firebase-conf.json")

	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	ds = datastore.Open(conf.DBconf)

	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/api/:api/organizations/:id/tickets.json", GetOrganizationTickets)
	//router.OPTIONS("/api/:api/organizations/:id/tickets.json", ConfigCORS)
	router.GET("/api/:api/organizations.json", GetOrganizations)
	//router.OPTIONS("/api/:api/organizations.json",ConfigCORS)
	router.OPTIONS("/*opts", ConfigCORS)
	log.Fatal(fasthttp.ListenAndServeTLS(":8080", "exclude/server.crt", "exclude/server.key", router.Handler))
}
