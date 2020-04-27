package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MrBTTF/gophercises/task/cmd"
	"github.com/MrBTTF/gophercises/task/db"
)

func main() {
	db, err := db.New("db/mytasks.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	ctx := context.WithValue(context.Background(), cmd.CtxDB, db)
	if err := cmd.Root.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
