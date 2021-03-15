//Numan client executable
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/footfish/numan"
	"github.com/footfish/numan/api/grpc"
	"github.com/footfish/numan/internal/app"
	"github.com/footfish/numan/internal/cmdcli"
	"github.com/gookit/color"
)

const (
	//DSN is path to sqlite file
	DSN = "./examples/numan-sqlite.db"
	//address = "localhost:50051"
)

func main() {
	cli := initCli()
	cli.Run()
}

//InitCli setup command configurations
func initCli() cmdcli.CommandConfigs {
	cli := cmdcli.NewCli()

	cmdDescription := "Adds a new number to the database. Number format is cc-ndc-sn"
	cmd := cli.NewCommand("add", add, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`) //mandatory params first.
	cmd.NewStringParameter("domain", true)
	cmd.NewStringParameter("carrier", true)

	cmdDescription = "Lists number db entries matching a number search. Number format is cc-ndc-sn, partial numbers are accepted "
	cmd = cli.NewCommand("list", list, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^([1-9]\d{0,2}\-[01]\d{0,4}\-\d{0,13})|([1-9]\d{0,2}\-[01]\d{0,4})$`)
	cmd.NewStringParameter("domain", false)

	cmdDescription = "Lists available numbers in db entries matching a number search. Number format is cc-ndc-sn, partial numbers are accepted "
	cmd = cli.NewCommand("list_free", listFree, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^([1-9]\d{0,2}\-[01]\d{0,4}\-\d{0,13})|([1-9]\d{0,2}\-[01]\d{0,4})$`)
	cmd.NewStringParameter("domain", false)

	cmdDescription = "Lists numbers for a user"
	cmd = cli.NewCommand("list_user", listUser, cmdDescription)
	cmd.NewIntParameter("uid", true)

	cmdDescription = "Views all details and history for number entries matching a number search. Number format is cc-ndc-sn, partial numbers are accepted  "
	cmd = cli.NewCommand("view", view, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{0,4}\-\d{1,13}$`)

	cmdDescription = "Deletes a number permentantly (history retained)"
	cmd = cli.NewCommand("delete", delete, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)

	cmdDescription = "Reserves a number for a user for a number of minutes"
	cmd = cli.NewCommand("reserve", reserve, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewIntParameter("uid", true)
	cmd.NewIntParameter("minutes", true).SetRegexp("^[0-9]{1,2}$")

	cmdDescription = "Sets a porting out date"
	cmd = cli.NewCommand("portout", portout, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewDateParameter("date", true)

	cmdDescription = "Sets a porting in date"
	cmd = cli.NewCommand("portin", portin, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewDateParameter("date", true)

	cmdDescription = "Allocates a number to a user"
	cmd = cli.NewCommand("allocate", allocate, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewIntParameter("uid", true)

	cmdDescription = "De-allocates a number from a user"
	cmd = cli.NewCommand("deallocate", deallocate, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)

	cmdDescription = "Provides a summary of number database"
	cmd = cli.NewCommand("summary", summary, cmdDescription)

	return cli
}

//newNuman instantiates a numan.API using either gRPC OR local database connection.
func newNuman() (nu numan.API) {
	address := os.Getenv("RPC_ADDRESS")
	if os.Getenv("RPC_ADDRESS") == "" {
		//fmt.Println("local db connection")
		return app.NewNumanService(DSN)
	}
	//fmt.Println("GRPC connection")
	return grpc.NewNumanClientAdapter(address)
}

//add <phonenumber> <domain> <carrier>
func add(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()

	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	newNumber := numan.Number{E164: numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2],
	},
		Domain:  p["domain"].(string),
		Carrier: p["carrier"].(string)}

	if err := nu.Add(&newNumber); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	}
	color.Info.Println("Success")
}

//list <phonenumber>
func list(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()

	var filter numan.NumberFilter
	splitNumber := strings.Split(p["phonenumber"].(string), "-")

	if len(splitNumber) == 2 {
		filter = numan.NumberFilter{E164: numan.E164{
			Cc:  splitNumber[0],
			Ndc: splitNumber[1]}}
	} else {
		filter = numan.NumberFilter{E164: numan.E164{
			Cc:  splitNumber[0],
			Ndc: splitNumber[1],
			Sn:  splitNumber[2]}}
	}
	if domain, ok := p["domain"].(string); ok {
		filter.Domain = domain
	}

	if numberList, err := nu.List(&filter); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.Info.Printf("Found %v\n", len(numberList))
		for _, n := range numberList {
			fmt.Println(colorize(fmt.Sprintf("%+v", n)))
		}
	}
}

//portout <phonenumber> <date>
func portout(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()

	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}
	portDate := p["date"].(time.Time).Unix()

	if err := nu.Portout(&number, &portDate); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.White.Print("Port date set")
	}
}

//portin <phonenumber> <date>
func portin(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()

	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}
	portDate := p["date"].(time.Time).Unix()

	if err := nu.Portin(&number, &portDate); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.White.Print("Port date set")
	}
}

//list_free <phonenumber>
func listFree(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()

	filter := numan.NumberFilter{State: 1} //1 = free
	splitNumber := strings.Split(p["phonenumber"].(string), "-")

	if len(splitNumber) == 2 {
		filter = numan.NumberFilter{E164: numan.E164{
			Cc:  splitNumber[0],
			Ndc: splitNumber[1]}}
	} else {
		filter = numan.NumberFilter{E164: numan.E164{
			Cc:  splitNumber[0],
			Ndc: splitNumber[1],
			Sn:  splitNumber[2]}}
	}

	if domain, ok := p["domain"].(string); ok {
		filter.Domain = domain
	}

	if numberList, err := nu.List(&filter); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.Info.Printf("Found %v\n", len(numberList))
		for _, n := range numberList {
			fmt.Println(colorize(fmt.Sprintf("%+v", n)))
		}
	}
}

//view <phonenumber>
func view(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()

	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}

	if numberDetails, err := nu.View(&number); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.White.Print(numberDetails)
	}
}

//summary
func summary(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()
	if summary, err := nu.Summary(); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.White.Print(summary)
	}
}

//delete <phonenumber>
func delete(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()

	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}

	if err := nu.Delete(&number); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.White.Print("Deleted")
	}
}

//reserve <phonenumber> <userid> <minutes>
func reserve(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}
	userID := p["uid"].(int64)
	untilTS := time.Now().Unix() + 60*p["minutes"].(int64)

	if err := nu.Reserve(&number, &userID, &untilTS); err != nil {
		color.Warn.Println(err)
		if numberDetails, err := nu.View(&number); err != nil {
			color.Warn.Println(err)
			os.Exit(1)
		} else {
			color.White.Print(numberDetails)
		}
		os.Exit(1)
	} else {
		color.Info.Println("Reserved")
	}
}

//allocate <phonenumber> <uid>
func allocate(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}
	userID := p["uid"].(int64)

	if err := nu.Allocate(&number, &userID); err != nil {
		color.Warn.Println(err)
		if numberDetails, err := nu.View(&number); err != nil {
			color.Warn.Println(err)
			os.Exit(1)
		} else {
			color.White.Print(numberDetails)
		}
		os.Exit(1)
	} else {
		color.Info.Println("Allocated")
	}
}

//deallocate <phonenumber>
func deallocate(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}

	if err := nu.DeAllocate(&number); err != nil {
		color.Warn.Println(err)
		if numberDetails, err := nu.View(&number); err != nil {
			color.Warn.Println(err)
			os.Exit(1)
		} else {
			color.White.Print(numberDetails)
		}
		os.Exit(1)
	} else {
		color.Info.Println("Deallocated")
	}
}

//	list_user <uid>
func listUser(p cmdcli.RxParameters) {
	nu := newNuman()
	defer nu.Close()
	userID := p["uid"].(int64)

	if numberList, err := nu.ListUserID(userID); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.Info.Printf("Found %v\n", len(numberList))
		for _, n := range numberList {
			fmt.Println(colorize(fmt.Sprintf("%+v", n)))
		}
	}
}

//colorize adds a bit of colour to formated strings (struct)
func colorize(s string) string {
	space := fmt.Sprintf(" "+color.SettingTpl, color.Cyan)
	colon := fmt.Sprintf(":"+color.SettingTpl, color.White)
	curley := fmt.Sprintf("{"+color.SettingTpl, color.Cyan)
	s = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, ":", colon), "{", curley), " ", space)
	return s
}
