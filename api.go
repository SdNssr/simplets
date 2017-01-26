package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
	"github.com/wcharczuk/go-chart"
	"log"
	"net/http"
)

func (env ServerEnv) ApiGetSeries(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	series, err := env.GetSeries()

	if err != nil {
		http.Error(w, fmt.Errorf("getting series: %s", err).Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(series)
}

func (env ServerEnv) ApiGetDataPoints(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	duration := r.URL.Query().Get("duration")

	datapoints, err := env.GetDataPoints(name, duration)

	if err != nil {
		http.Error(w, fmt.Errorf("getting datapoints: %s", err).Error(), 400)
		return
	}

	w.Header().Set("Content-Type", "image/png")

	graph := chart.Chart{
		Width:  2048,
		Height: 800,
		DPI:    192,
		XAxis: chart.XAxis{
			Style:          chart.StyleShow(),
			ValueFormatter: chart.TimeMinuteValueFormatter,
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: datapoints.x,
				YValues: datapoints.y,
			},
		},
	}

	err = graph.Render(chart.PNG, w)

	if err != nil {
		http.Error(w, fmt.Errorf("render graph: %s", err).Error(), 500)
	}
}

func (env ServerEnv) ApiAddDataPoint(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var dataPoint DataPoint

	if r.Body == nil {
		http.Error(w, "please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&dataPoint)

	if err != nil {
		http.Error(w, fmt.Errorf("decode json: %s", err).Error(), 400)
		return
	}

	name := p.ByName("name")
	value := dataPoint.Value

	err = env.AddDataPoint(name, value)

	if err != nil {
		http.Error(w, fmt.Errorf("add data point: %s", err).Error(), 400)
		return
	}

	fmt.Fprintf(w, "Data point added: %s: %f\n", name, value)
}

func ServeAPI(env ServerEnv) {
	router := httprouter.New()
	router.GET("/api/series", env.ApiGetSeries)
	router.GET("/api/series/:name", env.ApiGetDataPoints)
	router.POST("/api/series/:name", env.ApiAddDataPoint)

	n := negroni.Classic()
	n.UseHandler(router)

	log.Fatal(http.ListenAndServe(":8080", n))
}
