// Package cmdcli provides a simple command line.
// cmdcli parsers the first argument as the 'command' and subsequent arguments as 'parameters'.
// The format is 'command <param1>.. [param2]' (syntax <>mandatory []optional)
//
// Each command is linked to a command handling function.
// The string arguments are parsed/checked, then passed to the handling function as a map of interface{} (type converted) parameters (RxParameters)
//
// Usage:
// func main() {
// 	cli := cmdcli.NewCli()
// 	// Add command 'animal <type> <age> [name]'
// 	cmdDescription := "This command is used for something"
// 	cmd := cli.NewCommand("", animalHandler, cmdDescription)
// 	cmd.NewStringParameter("type", true).SetRegexp("^dog$|^cat$") //mandatory params first.
// 	cmd.NewIntParameter("age", true)
// 	cmd.NewStringParameter("name", false)
//
//  cli.Run()
//  }
// func animalHandler(p cmdcli.RxParameters) {
//     fmt.Println(p["animal"].(string),p["age"].(int))
// }
package cmdcli

import (
	"errors"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/gookit/color"
)

const (
	regexpParam     = "^[a-z][a-z0-9_]{0,20}$"
	regexpCommand   = "^[a-z][a-z0-9_]{0,20}$"
	regexpAnyString = ".*"
	regexpAnyInt    = "^[0-9]{1,10}$"
	regexpAnyFloat  = `^[0-9]{1,10}\.?[0-9]{0,10}$`
	regexpAnyDate   = `^\d{1,2}\/\d{1,2}\/2\d{3}$`
	typeString      = "string"
	typeInt         = "int"
	typeFloat       = "float"
	typeDate        = "date"
)

//RxParameters holds map of parsed parameters returned by cli to handler function
type RxParameters map[string]interface{}

//CommandConfigs holds configuration for all cli commands.
type CommandConfigs map[string]CommandConfig

//ParameterConfig stores cli commands parameter configurations
type ParameterConfig struct {
	fieldName string //name of the parameter
	fieldType string //convert string arg to parameter of variableType
	regexp    string //regular expression to validate string arg
	required  bool   //true if mandatory parameter
}

//CommandConfig stores cli commands configurations
type CommandConfig struct {
	param   map[int]*ParameterConfig //parameter configurations
	help    string
	handler func(RxParameters) //function which implements command. Passed a map of typed & parsed parameters.
}

type HandlerFunc func(RxParameters)

//NewCli makes the command configurations (1st function to call)
func NewCli() CommandConfigs {
	cli := make(CommandConfigs)
	return cli
}

//Run will parse the os.Args to a command parameter value map, then call the handler.
func (c CommandConfigs) Run() {
	var rxParams RxParameters
	switch len(os.Args) {
	case 1: //no command
		c.printHelp()
	case 2: //no parameters
		for _, p := range c.command(os.Args[1]).param {
			if p.required {
				color.Error.Println(p.fieldName + " required parameter")
				c.command(os.Args[1]).printUsage(os.Args[1])
				os.Exit(1)
			}
		}
		c.command(os.Args[1]).run(rxParams)
	default: //with parameters
		var err error
		if rxParams, err = c.command(os.Args[1]).readParameters(os.Args[2:]); err != nil {
			color.Error.Println(err)
			c.command(os.Args[1]).printUsage(os.Args[1])
			os.Exit(1)
		}
		c.command(os.Args[1]).run(rxParams)
	}
}

//NewCommand adds a new cli command.
func (c CommandConfigs) NewCommand(newCommand string, handlerFunc HandlerFunc, help string) CommandConfig {
	//validate param format
	if ok, err := regexp.MatchString(regexpCommand, newCommand); err != nil || !ok {
		panic(errors.New("Command " + newCommand + " not allowed. Regexp '" + regexpCommand + "'"))
	}
	//check unique
	if _, ok := c[newCommand]; ok {
		panic(errors.New("Command '" + newCommand + "' declared twice"))
	}
	c[newCommand] = CommandConfig{help: help, param: make(map[int]*ParameterConfig), handler: handlerFunc}
	return c[newCommand]
}

//command is getter for configCommand
func (c CommandConfigs) command(command string) CommandConfig {
	if _, ok := c[command]; !ok {
		color.Error.Println("Command", command, "does not exist")
		c.printHelp()
		os.Exit(1)
	}

	return c[command]
}

//printHelp prints usage & help for all commands
func (c CommandConfigs) printHelp() {
	color.Light.Println("\nGeneral Usage:-")
	color.Info.Println("\t" + os.Args[0] + " command <param1> [param2] [..]. ")
	color.Note.Println("\tSyntax: <mandatory> , [optional] ")
	color.Light.Println("\nSupported Commands:-")
	color.Info.Print("\t")

	keys := make([]string, 0, len(c)) // alphabetise command list
	for k := range c {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	first := true
	for _, k := range keys { // print comma separated command list
		if !first {
			color.Info.Print(", ")
		}
		color.Info.Print(k)
		first = false
	}
	color.Note.Println("\n\tUse command without parameters for detailed help.")

}

//parameter is getter for parameterConfig
func (c CommandConfig) parameter(findParam string) *ParameterConfig {
	var found *ParameterConfig
	for i, param := range c.param {
		if param.fieldName == findParam {
			found = c.param[i]
		}
	}
	return found
}

//readParameters will read and parse command line arguments from string array.
func (c CommandConfig) readParameters(args []string) (RxParameters, error) {
	rxParams := make(RxParameters)

	if len(c.param) == 0 {
		color.Note.Println("Info: Parameters ignored. The command does not need parameters.")
		return rxParams, nil //no args required
	}
	if len(c.param) < len(args) {
		color.Note.Println("Info: Some parameters ignored. The command takes ", len(c.param), " parameters")
	}
	for i, arg := range args {
		if i < len(c.param) {
			//reg-expr check.
			if ok, err := regexp.MatchString(c.param[i].regexp, arg); err != nil || !ok {
				return rxParams, errors.New("'" + arg + "' is invalid value for parameter " + c.param[i].fieldName + ". (Type '" + c.param[i].fieldType + "', Match expression '" + c.param[i].regexp + "')")
			}
			//convert to type.
			switch c.param[i].fieldType {
			case typeInt:
				var err error
				if rxParams[c.param[i].fieldName], err = strconv.ParseInt(arg, 10, 0); err != nil {
					return rxParams, errors.New("'" + arg + "' is invalid value for parameter " + c.param[i].fieldName + ". (Type '" + c.param[i].fieldType + "', Match expression '" + c.param[i].regexp + "')")
				}
			case typeString:
				rxParams[c.param[i].fieldName] = arg
			case typeFloat:
				var err error
				if rxParams[c.param[i].fieldName], err = strconv.ParseFloat(arg, 64); err != nil {
					return rxParams, errors.New("'" + arg + "' is invalid value for parameter " + c.param[i].fieldName + ". (Type '" + c.param[i].fieldType + "', Match expression '" + c.param[i].regexp + "')")
				}
			case typeDate:
				var err error
				if rxParams[c.param[i].fieldName], err = time.Parse("2/1/2006", arg); err != nil {
					return rxParams, errors.New("'" + arg + "' is invalid value for parameter " + c.param[i].fieldName + ". (Type '" + c.param[i].fieldType + "', Match expression '" + c.param[i].regexp + "')")
				}
			default:
				panic(errors.New("Command uses unsupported argument type '" + c.param[i].fieldType + "'"))
			}
		}
	}
	//check all mand. have been filled
	for _, p := range c.param {
		if p.required {
			if _, ok := rxParams[p.fieldName]; !ok {
				return rxParams, errors.New(p.fieldName + " is required parameter.")

			}
		}
	}
	return rxParams, nil
}

//run calls command handler with received parameter map
func (c CommandConfig) run(p RxParameters) {
	c.handler(p)
}

//NewStringParameter adds a new string parameter to command
func (c CommandConfig) NewStringParameter(newParam string, requiredness bool) *ParameterConfig {
	if err := c.validateParameter(newParam, requiredness); err != nil {
		panic(err)
	}
	i := len(c.param)
	c.param[i] = &ParameterConfig{fieldName: newParam, fieldType: typeString, regexp: regexpAnyString, required: requiredness}
	return c.param[i]
}

//NewDateParameter adds a new integer parameter to command
func (c CommandConfig) NewDateParameter(newParam string, requiredness bool) *ParameterConfig {
	if err := c.validateParameter(newParam, requiredness); err != nil {
		panic(err)
	}
	i := len(c.param)
	c.param[i] = &ParameterConfig{fieldName: newParam, fieldType: typeDate, regexp: regexpAnyDate, required: requiredness}
	return c.param[i]
}

//NewIntParameter adds a new integer parameter to command
func (c CommandConfig) NewIntParameter(newParam string, requiredness bool) *ParameterConfig {
	if err := c.validateParameter(newParam, requiredness); err != nil {
		panic(err)
	}
	i := len(c.param)
	c.param[i] = &ParameterConfig{fieldName: newParam, fieldType: typeInt, regexp: regexpAnyInt, required: requiredness}
	return c.param[i]
}

//NewFloatParameter adds a new floating point number parameter to command
func (c CommandConfig) NewFloatParameter(newParam string, requiredness bool) *ParameterConfig {
	if err := c.validateParameter(newParam, requiredness); err != nil {
		panic(err)
	}
	i := len(c.param)
	c.param[i] = &ParameterConfig{fieldName: newParam, fieldType: typeFloat, regexp: regexpAnyFloat, required: requiredness}
	return c.param[i]
}

//SetRegexp sets a regular expression for validation of parameter
func (p *ParameterConfig) SetRegexp(regexp string) *ParameterConfig {
	p.regexp = regexp
	return p
}

//validateParameter performs basic validation checks on new parameter
func (c CommandConfig) validateParameter(newParam string, requiredness bool) error {
	//validate param format
	if ok, err := regexp.MatchString(regexpParam, newParam); err != nil || !ok {
		return errors.New("Parameter '" + newParam + "' not allowed. Regexp '" + regexpParam + "'")
	}
	//check param is unique
	for _, param := range c.param {
		if param.fieldName == newParam {
			return errors.New("Parameter '" + newParam + "' already exists")
		}
	}
	//if required, then make sure previous parameter was required (obviously mandatory params must be before optional)
	if len(c.param) != 0 && requiredness {
		if c.param[len(c.param)-1].required == false {
			return errors.New("Parameter '" + newParam + "' is required, but previous was optional (must have mandatory params first)")
		}
	}

	return nil
}

//printUsage will print help and usage information for a command
func (c CommandConfig) printUsage(commandString string) {
	color.Light.Println("\nUsage:-")
	printString := "\t" + commandString
	for i := 0; i < len(c.param); i++ {
		if c.param[i].required {
			printString += " <" + c.param[i].fieldName + ">"
		} else {
			printString += " [" + c.param[i].fieldName + "] "
		}
	}
	color.Info.Println(printString)
	color.Note.Println("\t" + c.help)
}
