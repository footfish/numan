//Numan administrator executable
package main

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"log"
	"os"
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
	user numan.UserService
	ctx  context.Context //ctx ok here in structs as no scope issues. https://go.dev/blog/context-and-structs
	auth numan.User
}

var conf struct {
	Dsn           string
	ServerAddress string `envconfig:"optional"` //if ommitted works in standalone mode
	TlsCert       string `envconfig:"optional"` //if ommitted trusted Certificate Authority is needed
	TokenFile     string `envconfig:"default=.numa_auth"`
	User          string
	Password      string
}

func main() {
	var c client

	//Init conf from environmental vars
	godotenv.Load("numa.env")
	if err := envconfig.Init(&conf); err != nil {
		log.Fatalf("Failed to load required environmental variables for config: %v", err)
	}

	//Init context
	var cancel context.CancelFunc
	c.ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) //add client context
	defer cancel()

	//Init services
	if conf.ServerAddress == "" { //standalone servicelication with local db connection
		store := datastore.NewStore(conf.Dsn)
		defer store.Close()
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
		grpcClient := grpc.NewGrpcClient(c.ctx, conf.ServerAddress, creds)
		c.user = grpc.NewUserClientAdapter(grpcClient)
	}

	//Init authentication
	if err := c.setAuthToken(); err != nil { //load auth token from cached file or refresh
		color.Error.Println("Authentication error -", err)
		os.Exit(1)
	}
	c.ctx = context.WithValue(c.ctx, "token", c.auth.AccessToken) //add auth token to context

	//Run command line application
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

	cmdDescription := "Adds a new user to the database."
	cmd := cli.NewCommand("add", c.add, cmdDescription)
	cmd.NewStringParameter("username", true).SetRegexp(numan.PatternUser) //mandatory params first.
	cmd.NewStringParameter("password", true).SetRegexp(numan.PatternRawPassword)
	cmd.NewStringParameter("role", true).SetRegexp(`^(` + numan.RoleAdmin + `)|(` + numan.RoleUser + `)$`)

	cmdDescription = "Lists users. Will search partial usernames or list all."
	cmd = cli.NewCommand("list", c.list, cmdDescription)
	cmd.NewStringParameter("username", false)

	cmdDescription = "Delete a user"
	cmd = cli.NewCommand("delete", c.delete, cmdDescription)
	cmd.NewStringParameter("username", true).SetRegexp(numan.PatternUser) //mandatory params first.

	cmdDescription = "Sets a users password"
	cmd = cli.NewCommand("password", c.password, cmdDescription)
	cmd.NewStringParameter("username", true).SetRegexp(numan.PatternUser) //mandatory params first.
	cmd.NewStringParameter("password", true).SetRegexp(numan.PatternRawPassword)

	return cli
}

//add <username> <password> <role>
func (c *client) add(p cmdcli.RxParameters) {

	newUser := numan.User{
		Username: p["username"].(string),
		Role:     p["role"].(string),
		Password: p["password"].(string),
	}

	if err := c.user.AddUser(c.ctx, newUser); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	}

	color.Info.Println("Success, username '" + newUser.Username + "' added")
}

//list [username]
func (c *client) list(p cmdcli.RxParameters) {
	username, ok := p["username"].(string)
	if !ok {
		username = ""
	}

	userlist, err := c.user.ListUsers(c.ctx, username)
	if err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	}
	if len(userlist) == 0 {
		color.Warn.Println("None found")
		os.Exit(1)
	}
	printUserList(userlist)
}

//delete <username>
func (c *client) delete(p cmdcli.RxParameters) {
	username := p["username"].(string)

	if err := c.user.DeleteUser(c.ctx, username); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	}
	color.Info.Println("Deleted username '" + username + "'")
}

//password <username> <password>
func (c *client) password(p cmdcli.RxParameters) {
	username := p["username"].(string)
	password := p["password"].(string)
	if err := c.user.SetPassword(c.ctx, username, password); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	}
	color.Info.Println("New password set for username '" + username + "'")
}

//printUserList prints slice of numan.User as a table
func printUserList(userList []numan.User) {
	printer := tableprinter.New(os.Stdout)

	type tableRow struct {
		Username string `header:"Username"`
		Role     string `header:"Role"`
	}
	table := []tableRow{}

	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"

	for _, n := range userList {
		table = append(table, tableRow{
			Username: n.Username,
			Role:     n.Role,
		})
	}
	printer.Print(table)
}
