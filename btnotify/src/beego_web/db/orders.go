package db

import (
	"database/sql"
	//"database/sql/driver"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	mydb = &sql.DB{}
)

func init() {
	var err error
	mydb, err = sql.Open("mysql", "fabric:fabric@tcp(127.0.0.1:3306)/orders")
	if err != nil {
		fmt.Printf("open mysql fail, err=%v\r\n", err)
	}
}

func InsertOrder(sender, hash string, value int) error {
	_, err := mydb.Exec("INSERT INTO `orders` (`sender`, `hash`, `value`) VALUES (?,?,?)", sender, hash, value)
	if err != nil {
		fmt.Printf("insert sql err: %v\r\n", err)
		return err
	}
	return nil
}
