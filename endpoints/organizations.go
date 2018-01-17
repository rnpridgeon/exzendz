package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (e *endpoints) getOrganizations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(e.source.ExportOrganizations(0))
}

func (e *endpoints) getOrganizationTickets(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Println(r.Header)
	orgId, _ := strconv.ParseInt(p.ByName("id"), 10, 64)
	json.NewEncoder(w).Encode(e.source.ExportTickets(0, orgId))
}

func getOrganizationOptions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Under Contrstruction")
}

func (e *endpoints) ConfigureOrganizations() {
	e.router.OPTIONS("/api/organizations/*path", getOrganizationOptions)
	e.router.GET("/api/organizations", e.getOrganizations)
	e.router.GET("/api/organizations/:id", e.getOrganizationTickets)
}
