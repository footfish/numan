//Numan server executable
package main

import (
	"fmt"

	"github.com/footfish/numan/internal/app"
)

const (
	//DSN is path to sqlite file
	DSN = "./examples/numan-sqlite.db"
)

func main() {

	//TODO - build server
	//configure our core service
	nu := app.NewNumberService(DSN)
	defer nu.Close()
	fmt.Println(nu.Summary())
}
