package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

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

func parseIntRangeString(range_string string) (string, string, error) {
	elems := strings.Split(range_string, "-")
	if len(elems) != 2 {
		return "", "", errors.New("Invalid range parameter")
	}

	_, err := strconv.Atoi(elems[0])
	if err != nil {
		return "", "", err
	}
	_, err = strconv.Atoi(elems[1])
	if err != nil {
		return "", "", err
	}

	return elems[0], elems[1], nil
}

func parseFloatRangeString(range_string string) (string, string, error) {
	elems := strings.Split(range_string, "-")
	if len(elems) != 2 {
		return "", "", errors.New("Invalid range parameter")
	}

	_, err := strconv.ParseFloat(elems[0], 64)
	if err != nil {
		return "", "", err
	}
	_, err = strconv.ParseFloat(elems[1], 64)
	if err != nil {
		return "", "", err
	}

	return elems[0], elems[1], nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	db_host := os.Getenv("DB_HOST")
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	open_query := "postgres://" + db_user + ":" + db_password + "@" + db_host + "/" + db_name + "?sslmode=disable"

	db, err := sql.Open("pgx", open_query)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	fmt.Println("success to open")

	defer db.Close()

	query := "SELECT year FROM apartment_data"

	where_string := ""

	count_range, ok := request.QueryStringParameters["count"]
	if ok {
		min, max, err := parseIntRangeString(count_range)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		where_string += "count >= " + min + " AND count <= " + max
	}

	area_range, ok := request.QueryStringParameters["area"]
	if ok {
		min, max, err := parseFloatRangeString(area_range)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		if where_string != "" {
			where_string += " AND "
		}
		where_string += "area >= " + min + " AND area <= " + max
	}

	price_range, ok := request.QueryStringParameters["price"]
	if ok {
		min, max, err := parseIntRangeString(price_range)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}

		if where_string != "" {
			where_string += " AND "
		}
		where_string += "price >= " + min + " AND price <= " + max
	}

	if where_string != "" {
		query += " WHERE " + where_string
	}

	fmt.Println(query)
	rows, err := db.Query(query)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	var ret_list []int
	for rows.Next() {
		year := 0
		rows.Scan(&year)
		ret_list = append(ret_list, year)
	}

	ret, err := json.Marshal(ret_list)
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
