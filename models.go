package main

import "time"

type ticker struct {
        Name string `json:"name"`
        Title string `json:"title"`
        Price string `json:"price"`
}

type order struct {
        Ticker string `json:"ticker"`
        Price string `json:"price"`
}

type orderDB struct {
        Ticker string `json:"ticker"`
        Price float64 `json:"price"`
        Created time.Time `json:"created"`
}
