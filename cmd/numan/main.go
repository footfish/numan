//Numan client executable
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/footfish/numan"
	"github.com/footfish/numan/api/grpc"
	"github.com/footfish/numan/internal/app"
	"github.com/footfish/numan/internal/cmdcli"
	"github.com/footfish/numan/internal/datastore"
	"github.com/gookit/color"
	"google.golang.org/grpc/credentials"
)

const (
	//DSN is path to sqlite file
	dsn = "./examples/numan-sqlite.db"
	//address = "localhost:50051"
	certFile  = "./examples/server-cert.pem"
	tokenFile = ".numan_auth"
)

type client struct {
	numbering numan.NumberingService
	history   numan.HistoryService
	user      numan.UserService
	ctx       context.Context //TODO move out of struct
	auth      numan.User      //TODO do I need this now I have user.
}

func main() {
	var c client
	rpcAddress := os.Getenv("RPC_ADDRESS")
	if rpcAddress == "" { //standalone application with local db connection
		store := datastore.NewStore(dsn)
		defer store.Close()
		c.numbering = app.NewNumberingService(store)
		c.history = app.NewHistoryService(store)
		c.user = app.NewUserService(store)
	} else { //via gRPC
		creds := credentials.NewTLS(&tls.Config{})
		grpcClient := grpc.NewGrpcClient(rpcAddress, creds)
		c.numbering = grpc.NewNumberingClientAdapter(grpcClient)
		//		c.history = grpc.NewHistoryClientAdapter(rpcAddress, creds)
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
	if fileData, err := ioutil.ReadFile(tokenFile); err == nil {
		c.auth.AccessToken = string(fileData)
		//fmt.Println("Loaded token:", c.auth.AccessToken)
	}
	if c.auth.AuthRefreshRequired() {
		//need to fetch a token.
		//TODO - Testing, replace with config
		if c.auth, err = c.user.Auth(c.ctx, numan.USER, numan.PASSWORD); err != nil {
			return err
		}

		//Cache token
		if err := ioutil.WriteFile(tokenFile, []byte(c.auth.AccessToken), 0600); err != nil {
			color.Error.Println("Can't write to file -", tokenFile)
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

	cmdDescription = "Lists numbers for a user"
	cmd = cli.NewCommand("list_user", c.listUser, cmdDescription)
	cmd.NewIntParameter("uid", true)

	cmdDescription = "Views all details and history for number entries matching a number search. Number format is cc-ndc-sn, partial numbers are accepted  "
	cmd = cli.NewCommand("view", c.view, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{0,4}\-\d{1,13}$`)

	cmdDescription = "Deletes a number permentantly (history retained)"
	cmd = cli.NewCommand("delete", c.delete, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)

	cmdDescription = "Reserves a number for a user for a number of minutes"
	cmd = cli.NewCommand("reserve", c.reserve, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewIntParameter("uid", true)
	cmd.NewIntParameter("minutes", true).SetRegexp("^[0-9]{1,2}$")

	cmdDescription = "Sets a porting out date (dd/mm/yy)"
	cmd = cli.NewCommand("portout", c.portout, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewDateParameter("date", true)

	cmdDescription = "Sets a porting in date (dd/mm/yy)"
	cmd = cli.NewCommand("portin", c.portin, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewDateParameter("date", true)

	cmdDescription = "Allocates a number to a user"
	cmd = cli.NewCommand("allocate", c.allocate, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)
	cmd.NewIntParameter("uid", true)

	cmdDescription = "De-allocates a number from a user"
	cmd = cli.NewCommand("deallocate", c.deallocate, cmdDescription)
	cmd.NewStringParameter("phonenumber", true).SetRegexp(`^[1-9]\d{0,2}\-[01]\d{1,4}\-\d{5,13}$`)

	cmdDescription = "Provides a summary of number database"
	cmd = cli.NewCommand("summary", c.summary, cmdDescription)

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
		color.Info.Printf("Found %v\n", len(numberList))
		for _, n := range numberList {
			fmt.Println(colorize(fmt.Sprintf("%+v", n)))
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

	if numberList, err := c.numbering.List(c.ctx, &filter); err != nil {
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

//reserve <phonenumber> <userid> <minutes>
func (c *client) reserve(p cmdcli.RxParameters) {
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}
	userID := p["uid"].(int64)
	untilTS := time.Now().Unix() + 60*p["minutes"].(int64)

	if err := c.numbering.Reserve(c.ctx, &number, &userID, &untilTS); err != nil {
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

//allocate <phonenumber> <uid>
func (c *client) allocate(p cmdcli.RxParameters) {
	splitNumber := strings.Split(p["phonenumber"].(string), "-")
	number := numan.E164{
		Cc:  splitNumber[0],
		Ndc: splitNumber[1],
		Sn:  splitNumber[2]}
	userID := p["uid"].(int64)

	if err := c.numbering.Allocate(c.ctx, &number, &userID); err != nil {
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

	if err := c.numbering.DeAllocate(c.ctx, &number); err != nil {
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

//	list_user <uid>
func (c *client) listUser(p cmdcli.RxParameters) {
	userID := p["uid"].(int64)

	if numberList, err := c.numbering.ListUserID(c.ctx, userID); err != nil {
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
