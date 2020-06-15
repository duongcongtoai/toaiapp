package mysql

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func ConnectDB(url string) *gorm.DB {
	db, err := gorm.Open("mysql", url)
	if err != nil {
		fmt.Println(url)
		log.Fatalf("Can't connect to mysql server: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("%v captured - closing database connection\n", sig)
			db.Close()
			os.Exit(0)
		}
	}()
	return db
}
