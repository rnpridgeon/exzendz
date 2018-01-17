package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (e *endpoints) getTickets(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(e.source.ExportTickets(0, -1))
}

func getTicketOptions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Under Contrstruction")
}

func (e *endpoints) ConfigureTickets() {
	e.router.OPTIONS("/api/tickets/*path", getTicketOptions)
	e.router.GET("/api/tickets", e.getTickets)
	//e.router.GET("/tickets/:id", e.getTicket)
}
