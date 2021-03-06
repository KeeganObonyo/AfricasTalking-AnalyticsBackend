package main

import (
	"net/http"
	"time"
	"encoding/json"
    "github.com/gorilla/mux"
)

var path_details []string

//Server configuration
func main() {
	p("AFRICASTALKING REPO ANALYTICS", version(), "started at Adress:", config.Address, "\n", time.Now())

    router := mux.NewRouter()

	//The urls
	router.HandleFunc("/", Home).Name("Home")
	router.HandleFunc("/bar/graph/", GetGraph).Name("get_bar_graph")
	router.HandleFunc("/pie/chart/", GetGraph).Name("get_pie_chart")
    router.Walk(WalkFunc)

	server := &http.Server{
		Addr:           config.Address,
		Handler:        router,
		ReadTimeout:    time.Duration(config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
	return
}

//home
//lists urls at API homepage
func Home(writer http.ResponseWriter, request *http.Request) {
	{
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(writer)
		encoder.SetIndent(empty, tab)
		encoder.Encode(path_details)
	}

}
//Function for retrieving url details 
func WalkFunc(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
    url,_:=route.URLPath()
    path_details= append(path_details,url.Path)
    return nil
}
