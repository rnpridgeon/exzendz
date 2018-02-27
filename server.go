package main

import (
	"context"
	"log"
	"strconv"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/rnpridgeon/exzenzd/datastore"
	"github.com/rnpridgeon/utils/configuration"
	"os"
)

var (
	appCtx context.Context
	ds     *datastore.Datasource
)

type Config struct {
	//ZDconf *zendesk.ZendeskConfig `json:"zendesk"`
	Port     string            `json:"listener"`
	KeyFile  string            `json:"keyFile"`
	CertFile string            `json:"certFile"`
	DBconf   *datastore.Config `json:"database"`
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
		ctx.Response.Header.Add("Access-Control-Allow-Origin", Origin)
		ds.ExportOrganizations(ctx)
	}
}

func GetOrganizationTickets(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("Access-Control-Allow-Origin", Origin)
	if Authenticated(ctx) {
		orgid, _ := strconv.ParseInt(ctx.UserValue("id").(string), 10, 64)
		ds.ExportOrganizationTickets(orgid, ctx)
	}
}

func main() {
	var err error
	var conf Config

	if len(os.Args) != 3 {
		log.Fatalf("usage %s [exzenzd config] [authenticator config]\n", os.Args[0])
	}
	configuration.FromFile(os.Args[1])(&conf)
	InitAuthenticator(os.Args[2])

	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	ds = datastore.Open(conf.DBconf)

	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/api/:api/organizations/:id/tickets.json", GetOrganizationTickets)
	router.GET("/api/:api/organizations.json", GetOrganizations)
	router.OPTIONS("/*opts", ConfigCORS)
	log.Fatal(fasthttp.ListenAndServeTLS(conf.Port, conf.CertFile, conf.KeyFile, router.Handler))
}
