package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	db_host := os.Getenv("DB_HOST")
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	open_query := "postgres://" + db_user + ":" + db_password + "@" + db_host + "/" + db_name + "?sslmode=disable"
	fmt.Println(open_query)
	fmt.Println("host:" + db_host)
	fmt.Println("user:" + db_user)

	db, err := sql.Open("pgx", open_query)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	fmt.Println("success to open")

	defer db.Close()
	apartment_data_list := fetch_csv()

	current_rows, err := db.Query("SELECT year FROM apartment_data")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	var registered_years []int
	for current_rows.Next() {
		var year int
		current_rows.Scan(&year)
		registered_years = append(registered_years, year)
	}

	insert, err := db.Prepare("INSERT INTO apartment_data (year,count,area,price) VALUES ($1,$2,$3,$4)")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	for _, d := range apartment_data_list {
		if !contains(registered_years, d.Year) {
			insert.Exec(d.Year, d.Count, d.Area, d.Price)
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       "OK",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
