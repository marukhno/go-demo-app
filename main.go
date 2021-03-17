package main

import (
        "encoding/json"
        "io/ioutil"
        "log"
        "net/http"
        "strconv"
)

// A static map with tickers for GET /ticker?ticker=NAME method.
var tickers = map[string]ticker{
        "AAPL": {"AAPL", "Apple Inc.", "150.00"},
        "TSLA": {"TSLA", "Tesla Inc.", "800.50"},
        "IBM": {"IBM", "International Business Machines", "110.00"},
}

// pingHandler is an HTTP handler for GET /ping method. The method response with a static JSON.
func pingHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("{\"health\":\"ok\",\"getTicker\":\"/ticker?ticker=value\",\"getOrder\":\"/order?id=value\",\"postOrder\":\"/order -d {ticker:TSLA, price:1003.544}\"}"))
}

// orderHandler is an HTTP handler for GET /ticker method. It marshals JSON from a ticker structure and send it back to the requesting client.
func tickerHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
                ticker := r.URL.Query().Get("ticker")
                if v, ok := tickers[ticker]; ok {
                        resp, err := json.Marshal(v)
                        if err != nil {
                                w.Write([]byte("Failed to marshal JSON Ticker response\n"))
                                w.WriteHeader(http.StatusInternalServerError)
                                return
                        }
                        w.Header().Set("Content-Type", "application/json")
                        w.WriteHeader(http.StatusOK)
                        w.Write(resp)
                } else {
                        w.Header().Set("Content-Type", "application/json")
                        w.WriteHeader(http.StatusOK)
                        w.Write([]byte("{}"))
                }
        default:
                w.Header().Set("Allow", "GET")
                w.WriteHeader(http.StatusMethodNotAllowed)
        }
}

// orderHandler is an HTTP handler for GET /order?id=X and POST /order method. 
// For GET /order?id=X it select an order by its id from the database, marshal JSON response and send it to the client.
// For POST /order it expects JSON request like {"ticker":"AAPL", "price":"140.00"}, then unmarshal it and insert the data to the database. As a response it sends an ID of inserted order.
func orderHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
                id := r.URL.Query().Get("id")
                order, err := selectOrderId(connection(), id)
                if err != nil {
                        w.WriteHeader(http.StatusInternalServerError)
                        w.Write([]byte("Failed to select the order\n"))
                        w.Write([]byte(err.Error()))
                        return
                }
                resp, err := json.Marshal(order)
                if err != nil {
                        w.WriteHeader(http.StatusInternalServerError)
                        w.Write([]byte("Failed to marshal JSON Order response\n"))
                        return
                }
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                w.Write(resp)
        case "POST":
                body, _ := ioutil.ReadAll(r.Body)
                defer r.Body.Close()
                var order order
                err := json.Unmarshal(body, &order)
                if err != nil {
                        w.WriteHeader(http.StatusBadRequest)
                        w.Write([]byte("Failed to unmarshal JSON Order request\n"))
                        return
                }
                id, err := createOrder(connection(), order)
                if err != nil {
                        w.WriteHeader(http.StatusInternalServerError)
                        w.Write([]byte("Failed to place an order\n"))
                        w.Write([]byte(err.Error()))
                        return
                }
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusOK)
                response := `{"id":`+strconv.Itoa(id)+`}`
                w.Write([]byte(response))
        default:
                w.Header().Set("Allow", "GET, POST")
                w.WriteHeader(http.StatusMethodNotAllowed)
        }
}

// handleRequests registers handlers with our patterns.
func handleRequests() {
        http.HandleFunc("/ping", pingHandler)
        http.HandleFunc("/ticker", tickerHandler)
        http.HandleFunc("/order", orderHandler)
}

func main() {
        handleRequests()
        log.Println("starting server at :8080")
        http.ListenAndServe(":8080", nil)
}
