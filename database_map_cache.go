package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var mp map[int]string = make(map[int]string)

type data struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

func cache(id int) {
	db, err := sql.Open("mysql", "root:62011212@tcp(127.0.0.1:3306)/prodj")
	checkErr(err)

	_, in_cache := mp[id]

	if in_cache == true {
		// fmt.Println(mp[id])
		fmt.Println("-----------HIT----------")
	} else {
		fmt.Println("----------MISS----------")

		rows, err := db.Query("SELECT name, quantity_in_stock FROM products WHERE product_id = " + strconv.Itoa(id))
		checkErr(err)

		for rows.Next() {
			var name string
			var quantity int
			err = rows.Scan(&name, &quantity)

			result := data{Name: name, Quantity: quantity}
			byteArray, err := json.Marshal(result)
			checkErr(err)
			// fmt.Println(string(byteArray))

			mp[id] = string(byteArray)

		}
		// fmt.Println(mp)
		// fmt.Println("from data")
	}
}

func main() {
	for i := 0; i < 10; i++ {
		for j := 0; j < 2; j++ {
			start := time.Now()
			cache(i + 1)
			fmt.Printf("%v\n", time.Since(start))
		}
	}

}
