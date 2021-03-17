package main

import (
        "net/http"
        "net/http/httptest"
        "os"
        "testing"
)

func setEnv() {
        os.Setenv("DB_USERNAME", "docker")
        os.Setenv("DB_PASSWORD", "passw0rd")
        os.Setenv("DATABASE_URL", "localhost:5432/docker")
}

func TestServer(t *testing.T) {
        dbUsername := os.Getenv("DB_USERNAME")
        dbPassword := os.Getenv("DB_PASSWORD")
        dbUrl := os.Getenv("DATABASE_URL")
        setEnv()
        defer func(){
                os.Setenv("DB_USERNAME", dbUsername)
                os.Setenv("DB_PASSWORD", dbPassword)
                os.Setenv("DATABASE_URL", dbUrl)
        }()

        tt := []struct {
                name      string
                method    string
                target    string
                want      string
                statusCode int
        }{
                {
                        name:      "ping test",
                        method:    http.MethodGet,
                        target:    "/ping",
                        statusCode: http.StatusOK,
                },
                {
                        name:      "get order test",
                        method:    http.MethodGet,
                        target:    "/order",
                        statusCode: http.StatusOK,
                },
        }
        for _, tc := range tt {
                t.Run(tc.name, func(t *testing.T) {
                        responseRecorder := httptest.NewRecorder()

                        switch tc.target {
                        case "/ping":
                                request := httptest.NewRequest(tc.method, "/ping", nil)
                                pingHandler(responseRecorder, request)
                        case "/order":
                                request := httptest.NewRequest(tc.method, "/order?id=1", nil)
                                pingHandler(responseRecorder, request)
                        default:
                                t.Fatal("No method found for the test case")
                        }

                        if responseRecorder.Code != tc.statusCode {
                                t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
                        }

                })
        }
}
