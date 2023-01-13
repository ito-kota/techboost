package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// ToDo エラー処理
func fetch_csv() []ApartmentData {
	url := "https://www.toshiseibi.metro.tokyo.lg.jp/seisaku/tochi/data/tochi_2021_A_hyo1_3_1.csv"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	r := csv.NewReader(transform.NewReader(response.Body, japanese.ShiftJIS.NewDecoder()))
	year := 1983

	header, err := r.Read()

	count_column := -1
	area_column := -1
	price_column := -1

	for i, d := range header {
		if d == "供給戸数／都(戸)" {
			count_column = i
		} else if d == "1戸当たり平均住戸専有面積／都(m2)" {
			area_column = i
		} else if d == "1戸当たり平均住戸価格／都(万円)" {
			price_column = i
		}
	}

	var apartment_data_list []ApartmentData

	for {
		records, err := r.Read()

		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatal(err)
			panic(err)
		}

		count, _ := strconv.Atoi(records[count_column])
		area, _ := strconv.ParseFloat(records[area_column], 64)
		price, _ := strconv.Atoi(records[price_column])

		data := ApartmentData{Year: year, Count: count, Area: area, Price: price}

		apartment_data_list = append(apartment_data_list, data)
		year += 1
	}

	return apartment_data_list
}
