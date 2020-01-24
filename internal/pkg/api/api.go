package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Server exposes RESTful API endpoints.
type Server struct {
	Router *httprouter.Router
	logger *log.Logger
}

// NewServer instantiates an HTTP server and builds a route table.
func NewServer() *Server {
	s := &Server{
		Router: httprouter.New(),
		logger: log.New(os.Stderr, "API: ", 0),
	}
	s.buildRouteTable()
	return s
}

func (s *Server) buildRouteTable() {
	s.Router.GET("/api/probe", s.handleProbe())
	s.Router.GET("/api/rate", s.handleRate())
}

// encodeJSON takes content of any type (v) and encodes to the writer (w) in
// JSON format.
func (s *Server) encodeJSON(w http.ResponseWriter, v interface{}, l *log.Logger) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
	}
}

// handleProbe acts as the Kubernetes liveness & readiness probes.
//
// In most scenarios, liveness & readiness probes should probably do different
// things. For example, a readiness probe might include making a call to a
// database to ensure that functionality works. However for this exercise
// simply being able to get a 200 from an endpoint should suffice for both
// probes.
func (s *Server) handleProbe() httprouter.Handle {
	type response struct {
		Message string `json:"message"`
	}
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		payload := response{Message: "Application healthy!"}
		s.encodeJSON(w, payload, s.logger)
	}
}

// Rate is the valuation of a currency on a given date. This valuation compares
// the currency to Euros.
type Rate struct {
	Date     string  `json:"date"`
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

// handleRate fetches the latest value of £1 and $1 in €s. It also makes a
// recommendation whether or not it's worthwhile to exchange a currency at the
// current time based on historical data from the previous week.
func (s *Server) handleRate() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		params := r.URL.Query()
		s.logger.Printf("Calling exchangeratesapi.io/history\n")
		res, err := http.Get(buildURI(params.Get("currency")))
		if err != nil {
			Error(w, err, http.StatusServiceUnavailable)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			Error(w, err, http.StatusUnprocessableEntity)
		}

		var payload map[string]interface{}
		err = json.Unmarshal(body, &payload)
		if err != nil {
			Error(w, err, http.StatusUnprocessableEntity)
		}

		rates, err := getRates(payload)
		if err != nil {
			Error(w, err, http.StatusBadRequest)
		}

		s.encodeJSON(w, rates, s.logger)
	}
}

func buildURI(currency string) string {
	date := "2006-01-02"
	today := time.Now().Format(date)
	lastWeek := time.Now().AddDate(0, 0, -7).Format(date)
	uri := fmt.Sprintf("https://api.exchangeratesapi.io/history?start_at=%s&end_at=%s&symbols=%s", lastWeek, today, currency)
	return uri
}

func getRates(payload map[string]interface{}) ([]Rate, error) {
	var result []Rate

	rates, ok := payload["rates"].(map[string]interface{})
	if !ok {
		return []Rate{}, fmt.Errorf(`unable to parse "rates" from payload`)
	}

	// Since "date" is essentially a dynamic key, it make life easier later on
	// if the exchangeratesapi.io response is parsed into some more consistent
	// structure.
	for date, rate := range rates {
		r := Rate{}
		rate, ok := rate.(map[string]interface{})
		if !ok {
			return []Rate{}, fmt.Errorf("unable to parse exchange rate information")
		}
		for currency, value := range rate {
			r.Date = date
			r.Currency = currency
			r.Value = value.(float64)
			result = append(result, r)
		}
	}

	return result, nil
}
