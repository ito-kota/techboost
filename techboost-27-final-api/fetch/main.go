package main

import (
	"database/sql"
	"encoding/json"
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

type ApartmentData struct {
	Year  int     `json:"year"`
	Count int     `json:"count"`
	Area  float64 `json:"area"`
	Price int     `json:"price"`
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

	year, ok := request.QueryStringParameters["year"]
	query := "SELECT year,count,area,price FROM apartment_data"
	if ok {
		query += " WHERE year=" + year
	}

	rows, err := db.Query(query)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	ret_map := map[int]ApartmentData{}

	for rows.Next() {
		d := ApartmentData{}
		rows.Scan(&d.Year, &d.Count, &d.Area, &d.Price)
		ret_map[d.Year] = d
	}

	ret, err := json.Marshal(ret_map)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       string(ret),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
