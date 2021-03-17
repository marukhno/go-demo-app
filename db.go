package main

import (
        "context"
        "fmt"
        "github.com/jackc/pgx/v4"
        "log"
        "math"
        "os"
        "strconv"
        "time"
)

// connection creates a connection to the database using credentials from the environemnt variables.
func connection() *pgx.Conn {
        username := os.Getenv("DB_USERNAME")
        password := os.Getenv("DB_PASSWORD")
        dbUrl := "postgres://" + username +":" + password + "@" + os.Getenv("DATABASE_URL")
        conn, err := pgx.Connect(context.Background(), dbUrl)
        if err != nil {
                fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
                os.Exit(1)
        }
        return conn
}

// createOrder has a logic to insert our order into orders table.
func createOrder(conn *pgx.Conn, order order) (int, error) {
        defer func() {
                conn.Close(context.Background())
        }()
        sql := "insert into orders(ticker, price, created) values($1, $2, $3) RETURNING id;"
        ticker := order.Ticker
        price, err := strconv.ParseFloat(order.Price, 64)
        if err != nil {
                return -1, err
        }
        price = math.Round(price*100)/100
        currentTime := time.Now()
        created := currentTime.Format("2006-01-02 15:04:05")
        row := conn.QueryRow(context.Background(), sql, ticker, price, created)
        var id int
        err = row.Scan(&id)
        if err != nil {
                return -1, err
        }
        log.Printf("Inserted order for ticker=%s with id=%d\n", ticker, id)
        return id, nil
}

// selectOrderId has a logic to select an order by its id.
func selectOrderId(conn *pgx.Conn, id string) (orderDB, error) {
        defer func(){
                conn.Close(context.Background())
        }()

        var ticker string
        var price float64
        var created time.Time

        err := conn.QueryRow(context.Background(), "select ticker, price, created from orders where id=$1", id).Scan(&ticker, &price, &created)
        if err != nil {
                log.Println(err)
                return orderDB{}, err
        }

        order := orderDB{Price: price, Ticker: ticker, Created: created}
        return order, nil
}
