package main

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mrbttf/gophercises/phone"
	phonedb "github.com/mrbttf/gophercises/phone/db"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "phone_go"
	password = "phone_go"
	dbname   = "phone"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password)

	must(phonedb.Reset("postgres", psqlInfo, dbname))

	psqlInfo += " dbname=" + dbname
	must(phonedb.Migrate("postgres", psqlInfo))

	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	err = db.Seed()
	must(err)

	phones, err := db.AllPhones()
	must(err)

	for _, p := range phones {
		fmt.Printf("Working on %+v\n", p)
		number := phone.Normalize(p.Number)
		if number == p.Number {
			fmt.Println("No changes required...", number)
			continue
		}

		fmt.Println("Updating or removing...", number)
		existing, err := db.FindPhone(number)
		must(err)
		if existing != nil {
			must(db.DeletePhone(p.Id))
			continue
		}

		p.Number = number
		must(db.UpdatePhone(&p))

	}

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
