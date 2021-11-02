//Numan client executable
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/footfish/numan"
	"github.com/footfish/numan/api/grpc"
	"github.com/footfish/numan/internal/cmdcli"
	"github.com/footfish/numan/internal/service"
	"github.com/footfish/numan/internal/service/datastore"
	"github.com/gookit/color"
	"github.com/joho/godotenv"
	"github.com/lensesio/tableprinter"
	"github.com/vrischmann/envconfig"
	"google.golang.org/grpc/credentials"
)

type client struct {
	numbering numan.NumberingService
	history   numan.HistoryService
	user      numan.UserService
	ctx       context.Context //TODO move out of struct
	auth      numan.User      //TODO do I need this now I have user.
}

var conf struct {
	Dsn           string
	ServerAddress string `envconfig:"optional"` //if ommitted works in standalone mode
	TlsCert       string `envconfig:"optional"` //if ommitted trusted Certificate Authority is needed
	TokenFile     string `envconfig:"default=.numan_auth"`
	User          string
	Password      string
}

func main() {
	//Init conf from environmental vars
	godotenv.Load("numan.env")
	if err := envconfig.Init(&conf); err != nil {
		log.Fatalf("Failed to load required environmental variables for config: %v", err)
	}

	var c client
	if conf.ServerAddress == "" { //standalone servicelication with local db connection
		store := datastore.NewStore(conf.Dsn)
		defer store.Close()
		c.numbering = service.NewNumberingService(store)
		c.history = service.NewHistoryService(store)
		c.user = service.NewUserService(store)
	} else { //via gRPC
		var creds credentials.TransportCredentials
		if conf.TlsCert == "" { //Using trusted CA, no need to load client cert
			creds = credentials.NewTLS(&tls.Config{})
		} else { //use self-signed cert
			var err error
			creds, err = credentials.NewClientTLSFromFile(conf.TlsCert, "")
			if err != nil {
				log.Fatalf("cert load error: %s", err)
			}
		}
		grpcClient := grpc.NewGrpcClient(conf.ServerAddress, creds)
		c.numbering = grpc.NewNumberingClientAdapter(grpcClient)
		c.history = grpc.NewHistoryClientAdapter(grpcClient)
		c.user = grpc.NewUserClientAdapter(grpcClient)
	}

	var cancel context.CancelFunc
	c.ctx, cancel = context.WithTimeout(context.Background(), time.Second) //add client context
	defer cancel()

	//Get authentication token in context
	//store token in context.
	if err := c.setAuthToken(); err != nil {
		color.Error.Println("Authentication error -", err)
		os.Exit(1)
	}
	c.ctx = context.WithValue(c.ctx, "token", c.auth.AccessToken) //add auth token to context
	c.initCli().Run()
}

func (c *client) setAuthToken() (err error) {
	//load file token
	if fileData, err := ioutil.ReadFile(conf.TokenFile); err == nil {
		c.auth.AccessToken = string(fileData)
		//fmt.Println("Loaded token:", c.auth.AccessToken)
	}
	if c.auth.AuthRefreshRequired() {
		//need to fetch a token.
		if c.auth, err = c.user.Auth(c.ctx, conf.User, conf.Password); err != nil {
			return err
		}

		//Cache token
		if err := ioutil.WriteFile(conf.TokenFile, []byte(c.auth.AccessToken), 0600); err != nil {
			color.Error.Println("Can't write to file -", conf.TokenFile)
			os.Exit(1)
		}
		return nil
	}
	return nil
}

//InitCli setup command configurations
func (c *client) initCli() cmdcli.CommandConfigs {
	cli := cmdcli.NewCli()

	cmdDescription := "Adds a new number to the database. Number format is cc-ndc-sn"
	cmd := cli.NewCommand("add", c.add, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`) //mandatory params first.
	cmd.NewStringParameter("domain", true)
	cmd.NewStringParameter("carrier", true)

	cmdDescription = "Lists number db entries matching a number search. Number format is cc-ndc-sn, partial numbers are accepted "
	cmd = cli.NewCommand("list", c.list, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^([1-9]\d{0,2}\-[01]\d{0,4}\-\d{0,13})|([1-9]\d{0,2}\-[01]\d{0,4})$`)
	cmd.NewStringParameter("domain", false)

	cmdDescription = "Lists available numbers in db entries matching a number search. Number format is cc-ndc-sn, partial numbers are accepted "
	cmd = cli.NewCommand("list_free", c.listFree, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^([1-9]\d{0,2}\-[01]\d{0,4}\-\d{0,13})|([1-9]\d{0,2}\-[01]\d{0,4})$`)
	cmd.NewStringParameter("domain", false)

	cmdDescription = "Lists numbers attached to owner"
	cmd = cli.NewCommand("list_owner", c.listOwner, cmdDescription)
	cmd.NewIntParameter("oid", true)

	cmdDescription = "Views all details and history for number entries matching a number search. Number format is cc-ndc-sn, partial numbers are accepted  "
	cmd = cli.NewCommand("view", c.view, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{0,4}\-\d{1,13}$`)

	cmdDescription = "Deletes a number permentantly (history retained)"
	cmd = cli.NewCommand("delete", c.delete, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)

	cmdDescription = "Reserves a number for an owner for a number of minutes"
	cmd = cli.NewCommand("reserve", c.reserve, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewIntParameter("oid", true)
	cmd.NewIntParameter("minutes", true).SetRegexp("^[0-9]{1,2}$")

	cmdDescription = "Sets a porting out date (dd/mm/yyyy)"
	cmd = cli.NewCommand("portout", c.portout, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewDateParameter("date", true)

	cmdDescription = "Sets a porting in date (dd/mm/yyyy)"
	cmd = cli.NewCommand("portin", c.portin, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewDateParameter("date", true)

	cmdDescription = "Allocates a number to an owner"
	cmd = cli.NewCommand("allocate", c.allocate, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewIntParameter("oid", true)

	cmdDescription = "De-allocates a number from an owner"
	cmd = cli.NewCommand("deallocate", c.deallocate, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewIntParameter("oid", true)

	cmdDescription = "Provides a summary of number database"
	cmd = cli.NewCommand("summary", c.summary, cmdDescription)

	cmdDescription = "Lists history log for an owner"
	cmd = cli.NewCommand("history_owner", c.listHistoryOwner, cmdDescription)
	cmd.NewIntParameter("oid", true)

	cmdDescription = "Lists history log for a number"
	cmd = cli.NewCommand("history", c.listHistoryNumber, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	return cli
}

//add <phonenumber> <domain> <carrier>
func (c *client) add(p cmdcli.RxParameters) {
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	newNumber := numan.Numbering{E164: numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2],
	},
		Domain:  p["domain"].(string),
		Carrier: p["carrier"].(string)}

	if err := c.numbering.Add(c.ctx, &newNumber); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	}
	color.Info.Println("Success")
}

//list <phonenumber>
func (c *client) list(p cmdcli.RxParameters) {
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

	if numberList, err := c.numbering.List(c.ctx, &filter); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		printNumberList(numberList)
		if len(numberList) == 1 { //print number history if there is only one result.
			if historyList, err := c.history.ListHistoryByNumber(c.ctx, numberList[0].E164); err != nil {
				color.Warn.Println(err)
				os.Exit(1)
			} else {
				printHistoryList(historyList)
			}
		}
	}

}

//portout <phonenumber> <date>
func (c *client) portout(p cmdcli.RxParameters) {
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}
	portDate := p["date"].(time.Time).Unix()

	if err := c.numbering.Portout(c.ctx, &number, &portDate); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.White.Print("Port date set")
	}
}

//portin <phonenumber> <date>
func (c *client) portin(p cmdcli.RxParameters) {
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}
	portDate := p["date"].(time.Time).Unix()

	if err := c.numbering.Portin(c.ctx, &number, &portDate); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.White.Print("Port date set")
	}
}

//list_free <phonenumber>
func (c *client) listFree(p cmdcli.RxParameters) {
	var filter numan.NumberFilter
	splitNumber := strings.Split(p["phonenumber"].(string), "-")

	if len(splitNumber) == 2 {
		filter = numan.NumberFilter{
			E164: numan.E164{
				Cc:  splitNumber[0],
				Ndc: splitNumber[1]},
			State: 1, //free
		}
	} else {
		filter = numan.NumberFilter{
			E164: numan.E164{
				Cc:  splitNumber[0],
				Ndc: splitNumber[1],
				Sn:  splitNumber[2]},
			State: 1, //free
		}
	}

	if domain, ok := p["domain"].(string); ok {
		filter.Domain = domain
	}

	if numberList, err := c.numbering.List(c.ctx, &filter); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		printNumberList(numberList)
	}
}

//view <phonenumber>
func (c *client) view(p cmdcli.RxParameters) {
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}

	if numberDetails, err := c.numbering.View(c.ctx, &number); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.White.Print(numberDetails)
	}
}

//summary
func (c *client) summary(p cmdcli.RxParameters) {
	if summary, err := c.numbering.Summary(c.ctx); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.White.Print(summary)
	}
}

//delete <phonenumber>
func (c *client) delete(p cmdcli.RxParameters) {
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}

	if err := c.numbering.Delete(c.ctx, &number); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		color.White.Print("Deleted")
	}
}

//reserve <phonenumber> <oid> <minutes>
func (c *client) reserve(p cmdcli.RxParameters) {
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}
	ownerID := p["oid"].(int64)
	untilTS := time.Now().Unix() + 60*p["minutes"].(int64)

	if err := c.numbering.Reserve(c.ctx, &number, &ownerID, &untilTS); err != nil {
		color.Warn.Println(err)
		if numberDetails, err := c.numbering.View(c.ctx, &number); err != nil {
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

//allocate <phonenumber> <oid>
func (c *client) allocate(p cmdcli.RxParameters) {
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}
	ownerID := p["oid"].(int64)

	if err := c.numbering.Allocate(c.ctx, &number, &ownerID); err != nil {
		color.Warn.Println(err)
		if numberDetails, err := c.numbering.View(c.ctx, &number); err != nil {
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
func (c *client) deallocate(p cmdcli.RxParameters) {
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}
	ownerID := p["oid"].(int64)

	if err := c.numbering.DeAllocate(c.ctx, &number, &ownerID); err != nil {
		color.Warn.Println(err)
		if numberDetails, err := c.numbering.View(c.ctx, &number); err != nil {
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

//	list_owner <oid>
func (c *client) listOwner(p cmdcli.RxParameters) {
	ownerID := p["oid"].(int64)

	if numberList, err := c.numbering.ListOwnerID(c.ctx, ownerID); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		printNumberList(numberList)
	}
}

//printNumberList prints slice of numan.Numbering as a table
func printNumberList(numberList []numan.Numbering) {
	printer := tableprinter.New(os.Stdout)

	type tableRow struct {
		ID          int64  `header:"ID"`
		Number      string `header:"Number"`
		Domain      string `header:"Domain"`
		Carrier     string `header:"Carrier"`
		OwnerID     int64  `header:"Owner"`
		Used        bool   `header:"Used"`
		Allocated   string `header:"Allocated"`
		Reserved    string `header:"Reserved"`
		DeAllocated string `header:"De-alloc'd"`
		PortedIn    string `header:"Port IN"`
		PortedOut   string `header:"Port OUT"`
	}
	table := []tableRow{}

	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"

	dateConv := func(unixTime int64) string {
		if unixTime == 0 {
			return "-"
		}
		return time.Unix(unixTime, 0).Format(numan.DATEPRINTFORMAT)
	}

	for _, n := range numberList {
		table = append(table, tableRow{ID: n.ID,
			OwnerID:     n.OwnerID,
			Domain:      n.Domain,
			Carrier:     n.Carrier,
			Number:      fmt.Sprintf("%v-%v-%v", n.E164.Cc, n.E164.Ndc, n.E164.Sn),
			Used:        n.Used,
			Allocated:   dateConv(n.Allocated),
			Reserved:    dateConv(n.Reserved),
			DeAllocated: dateConv(n.DeAllocated),
			PortedIn:    dateConv(n.PortedIn),
			PortedOut:   dateConv(n.PortedOut),
		})
	}
	printer.Print(table)
}

//	history_owner <oid>
func (c *client) listHistoryOwner(p cmdcli.RxParameters) {
	ownerID := p["oid"].(int64)

	if historyList, err := c.history.ListHistoryByOwnerID(c.ctx, ownerID); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		printHistoryList(historyList)
	}
}

//	history <phonenumber>
func (c *client) listHistoryNumber(p cmdcli.RxParameters) {
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}

	if historyList, err := c.history.ListHistoryByNumber(c.ctx, number); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	} else {
		printHistoryList(historyList)
	}
}

//printHistoryList prints slice of numan.History as a table
func printHistoryList(historyList []numan.History) {
	printer := tableprinter.New(os.Stdout)

	type tableRow struct {
		Timestamp string `header:"Timestamp"`
		Action    string `header:"Action"`
		Number    string `header:"Number"`
		OwnerID   int64  `header:"Owner"`
		Notes     string `header:"Notes"`
	}
	table := []tableRow{}

	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"

	dateConv := func(unixTime int64) string {
		if unixTime == 0 {
			return "-"
		}
		return time.Unix(unixTime, 0).Format(numan.TIMESTAMPPRINTFORMAT)
	}

	for _, n := range historyList {
		table = append(table, tableRow{
			Timestamp: dateConv(n.Timestamp),
			Action:    n.Action,
			Number:    fmt.Sprintf("%v-%v-%v", n.E164.Cc, n.E164.Ndc, n.E164.Sn),
			OwnerID:   n.OwnerID,
			Notes:     n.Notes,
		})
	}
	printer.Print(table)
}
