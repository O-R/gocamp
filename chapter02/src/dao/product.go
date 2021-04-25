package dao

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
)

type Model struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Product struct {
	Model
	Code  string
	Price uint
}

func GetProducts() error {
	db, err := sql.Open(sqlite.DriverName, "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	var (
		code  string
		price uint
	)
	var products []Product

	// rows, err := db.Query("select code, price from products where id = ?", 0)
	rows, err := db.Query("select code, price from products")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	//https://golang.org/pkg/database/sql/#pkg-variables
	//ErrNoRows is returned by Scan when QueryRow doesn't return a row. In such a case, QueryRow returns a placeholder *Row value that defers this error until a Scan.
	//从定义可知 ErrNoRows 在使用不当时才会发生，正确使用 Next 不会有这个 error
	for rows.Next() {
		err := rows.Scan(&code, &price)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, Product{Code: code, Price: price})
		log.Println(code, price)
	}

	fmt.Printf("products: %+v", products)
	return rows.Err()
}
