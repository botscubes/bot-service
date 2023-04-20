package pgsql

import (
	"context"
	"fmt"
	"os"
)

func (db *Db) GetTest() {
	var version string
	err := db.Pool.QueryRow(context.Background(), "select version()").Scan(&version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(version)
}
