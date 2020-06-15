package postgresql

import (
	"log"
	"os"
	"os/signal"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ConnectDB(url string) *gorm.DB {
	db, err := gorm.Open("postgres", url)
	if err != nil {
		log.Fatalf("Can't connect to postgresql server: %v", err)
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
