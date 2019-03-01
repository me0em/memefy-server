// Package db represents databases abstractions
package api

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

/*
type Product struct {
	gorm.Model
	Code string
	Price uint
}

func TestDB() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})
	fmt.Println("1")

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})
	fmt.Println("2")

	db.Create(&Product{Code: "L121ewf2", Price: 10020})
	fmt.Println("2")

	db.Create(&Product{Code: "fewrf", Price: 143000})
	fmt.Println("2")
}

func TestDB() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{})

	// Create
	db.Create(&User{UserID: "test-user-id", IDType: "Google oauth"})

	// Read
	//
	//var product Product
	//db.First(&product, 1) // find product with id 1
	//db.First(&product, "code = ?", "L1212") // find product with code l1212
	//
	//// Update - update product's price to 2000
	//db.Model(&product).Update("Price", 2000)
	//
	//// Delete - delete product
	//db.Delete(&product)
}

*/

func TestDB() {
	fmt.Printf("DB operation")
}