package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
)

var options struct {
	writeNum int
	readNum  int
}

func init() {
	flag.IntVar(&options.writeNum, "w", 1, "")
	flag.IntVar(&options.readNum, "r", 10, "")
}

func main() {
	flag.Parse()
	log.Printf("writeNum: %d, readNum: %d\n", options.writeNum, options.readNum)

	db, err := sql.Open("mysql", "user1:pass@tcp(moco-target-primary.default.svc:3306)/db1")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	write(db, options.writeNum)
	read(db, options.readNum)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh
	log.Print("Bye!")
}

func write(db *sql.DB, num int) {
	for i := 0; i < num; i++ {
		go func(index int) {
			log.Printf("write[%d]: start", index)
			for {
				_, err := db.Exec(`INSERT INTO db1.t1 (v1) VALUES (1)`)
				if err != nil {
					log.Printf("write[%d]: Error %v", index, err)
				} else {
					log.Printf("write[%d]: OK", index)
				}
			}
		}(i)
	}
}

func read(db *sql.DB, num int) {
	for i := 0; i < num; i++ {
		go func(index int) {
			log.Printf("read[%d]: start", index)
			for {
				row := db.QueryRow(`SELECT COUNT(*) from db1.t1 LIMIT 1000`)
				if err := row.Err(); err != nil {
					log.Printf("read[%d]: Error %v", index, err)
				} else {
					log.Printf("read[%d]: OK", index)
				}
			}
		}(i)
	}
}
