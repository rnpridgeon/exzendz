package endpoints

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rnpridgeon/zendb/models"
	"github.com/rnpridgeon/zendb/provider/mysql"
	"net/http"
)

// TODO: Figure out what was missing from the interface
type DataStore interface {
	ExportOrganizations(int64, func([]models.Organization) int64)
	ExportTickets(int64, orgId int64, fn func([]models.Ticket_Enhanced) (last int64))
}

type endpoints struct {
	source *mysql.MysqlProvider
	router *httprouter.Router
}

func Configure(router *httprouter.Router, source *mysql.MysqlProvider) {
	e := &endpoints{source, router}
	e.ConfigureOrganizations()
	e.ConfigureTickets()
	router.NotFound = http.FileServer(http.Dir("/Users/ryan/Desktop/react-corse-projects/digest/public"))
}
