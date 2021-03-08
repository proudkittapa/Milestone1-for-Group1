package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db    *sql.DB
	mutex sync.Mutex
)

func getQuantity(id int) int {
	row, err := db.Query("select * from products where product_id = " + strconv.Itoa(id))
	if err != nil {
		panic(err)
	}
	var q int
	for row.Next() {
		var id int
		var name string
		row.Scan(&id, &name, &q)
		fmt.Println("id: ", id, " product: ", name, " quantity: ", q)
	}
	return q
}
func decrement(quantity int, orderQuantity int, id int) {
	newQuantity := quantity - orderQuantity
	if newQuantity < 0 {
		return
	}
	fmt.Println(newQuantity)
	db.Query("update products set quantity_in_stock = ? where product_id = ? ", newQuantity, id)
}

func insert(user string, id int, q int) {
	db.Query("INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
}

func preorder(end chan int, user string, productId int, orderQuantity int) {
	// fmt.Printf("start\n")
	start := time.Now()
	mutex.Lock()
	quantity := getQuantity(productId)
	decrement(quantity, orderQuantity, productId)
	insert(user, productId, orderQuantity)
	mutex.Unlock()
	fmt.Printf("time: %v\n", time.Since(start))
	num, _ := strconv.Atoi(user)
	end <- num
	return
}
func main() {
	db, _ = sql.Open("mysql", "root:62011212@tcp(127.0.0.1:3306)/prodj")
	end := make(chan int)
	for i := 0; i < 5; i++ {
		go preorder(end, strconv.Itoa(i), 1, 5)
	}
	for i := range end {
		fmt.Println(i)
	}
}
