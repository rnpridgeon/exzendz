package datastore

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/allegro/bigcache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"

	"github.com/rnpridgeon/zendb/provider/mysql"
)

const (
	dsn = "%v:%s@tcp(%s:%d)/zendb?charset=utf8"
)

type Request struct {
	Resource string
	Id       string
	Filter   string
}

type Config struct {
	mysql.MysqlConfig
	Expiry int `json:"cache_expiry_min"`
}

type Datasource struct {
	dbClient   *sqlx.DB
	statements map[string]string
	cache      *bigcache.BigCache
}

func Open(conf *Config) *Datasource {
	db, err := sqlx.Open(conf.Type, fmt.Sprintf(dsn,
		conf.User, conf.Password, conf.Hostname, conf.Port))

	err = db.Ping()

	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}

	cacheConfig := bigcache.DefaultConfig(1 * time.Minute)
	cacheConfig.CleanWindow = time.Duration(conf.Expiry) * time.Minute
	cacheConfig.HardMaxCacheSize = 50

	cache, err := bigcache.NewBigCache(cacheConfig)

	if err != nil {
		log.Fatal("Failed to initialize cache:", err)
	}

	return &Datasource{
		dbClient:   db,
		statements: make(map[string]string),
		cache:      cache,
	}
}

func (ds *Datasource) ExportOrganizations(ctx *fasthttp.RequestCtx) {
	if !ds.get("Organizations", ctx) {
		var Organizations []OrganizationView

		fmt.Println(ds.dbClient.Select(&Organizations, "SELECT * FROM OrganizationView WHERE OrganizationView.Name not like '%deleted%' ORDER BY OrganizationView.id"))

		val, err := json.Marshal(Organizations)

		if err != nil {
			log.Println("ERROR: Failed to process request: ", err)
			ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)
			return
		}

		ctx.SetBody(val)
		ctx.SetStatusCode(fasthttp.StatusOK)

		ds.cache.Set("Organizations", val)
		log.Println("INFO: Added Resultset for Organizations to cache")
	}
}

func (ds *Datasource) ExportOrganizationTickets(organizationid int64, ctx *fasthttp.RequestCtx) {
	orgid := strconv.FormatInt(organizationid, 10)
	if !ds.get("OrganizationTickets:"+orgid, ctx) {
		var Tickets []TicketView

		fmt.Println(ds.dbClient.Select(&Tickets, ""+
			"SELECT TicketView.* "+
			"	FROM TicketView "+
			"	JOIN ticket ON TicketView.id = ticket.id "+
			" WHERE TicketView.organizationid = ? ORDER BY TicketView.id", orgid))

		val, err := json.Marshal(Tickets)

		if err != nil {
			log.Println("ERROR: Failed to process request: ", err)
			ctx.SetStatusCode(fasthttp.StatusUnprocessableEntity)
			return
		}

		ctx.SetBody(val)
		ctx.SetStatusCode(fasthttp.StatusOK)

		ds.cache.Set("OrganizationTickets:"+orgid, val)
		log.Printf("INFO: Added Resultset for OrganizationTickets:%s to cache\n", orgid)
	}
}

func (ds *Datasource) get(key string, ctx *fasthttp.RequestCtx) (OK bool) {
	if val, _ := ds.cache.Get(key); val != nil {
		log.Printf("Returning recordset for %s from cache\n", key)
		ctx.SetBody(val)
		return true
	}

	return false
}
