# numan - a phone number management tool 

This is an example Go gRPC microservice project.
It's a simple command line tool to keep track of phone number allocations/reservations (for a service provider). Data is stored in an sqlite db. 

```
General Usage:-
        numan command <param1> [param2] [..]. 
                Syntax: <mandatory> , [optional] 

Example:-
        numan view 353-01-12345111    # view details for phone number 353-01-12345111
```
It's a personal (learning) project (to replace an excel file). The main purpose is to:
- explore suitable Go project layout. 
- explore using gRPC (and perhaps extending to expose gRPC via REST ). 

The project loosely uses [DDD (domain driven design)](https://en.wikipedia.org/wiki/Domain-driven_design) and the  [standard Go project layout](https://github.com/golang-standards/project-layout). 

DDD's main principle is that classes should match the business domain. In this case the main _business objects_ being  a change 'history' log and the 'numbering' & 'user' records. 

For this project the _business objects_ (history, numbering, user) are presented as separate service interfaces (defined in the root folder; history.go, numbering.go, user.go), then service/object logic is 'layered' in the service directory (I guess you could call this Go's version of method overloading ). 

To demonstrate using an example from the project folders. 

- 'numbering' is a _business object_ implemented as a service
- /numbering.go <- this defines the _business object API_ ie. interface and structs  
- /internal/service/numbering.go <- this implements the numbering interface to the core service 
- /internal/service/datastore/numbering.go <- this implements the numbering storage layer

This pattern allows for: 
- Separation of concerns 
- Easily changing/adding new layers. For example swapping to a different storage mechanism. 
- Avoids circular dependencies (which go does not allow). 

## Installation
This installation uses Go modules, so should not be installed in your $GOPATH (of course you will need Go installed to compile)

```
$ git clone https://github.com/footfish/numan
$ cd numan
$ go install ./cmd/...
 
```
Two binaries are installed 
*  numan (client or standalone)
*  numand (server)

## Configuration 
Configuration is from environmental variables. 
These can be loaded from files or set on the terminal. 
Check ./examples folder for list (*.env files).

## Running Standalone Mode

numan can be run as standalone or as client/server(gRPC).
To run as standalone application make sure the env SERVER_ADDRESS is NOT configured.

Try the example: 
```
$ cd examples                   # contains sample config and db. 

$ vi numan.env                  # remove/comment line SERVER_ADDRESS

$ numan                         # prints the command help 
General Usage:-
        command <param1> [param2] [..]. 
                Syntax: <mandatory> , [optional]

# Sample usage (see /scripts folder for more examples)
$ numan summary                 # prints a summary of the example database 
$ numan list 353-01             # prints details of numbers starting 353-01
$ numan view 353-01-12345111    # view all details for number 353-01-12345111

```
### Running Client-Server Mode

To run in client-server mode you will need to use certificates. 
There are two approaches you can use. 
 1) Use a self-signed cert. This approach requires the cert to be loaded in the client (or switch off verification). 
 2) Use your own trusted CA with minica (preferred). This approach does NOT require the cert to be loaded in the client but it a little trickier to set up. 

 The example uses a self-signed cert. For help installing your own certs - see [/scripts/gen_certs.sh](./scripts/gen_certs.sh)

The client requires authentication. The username/password 

Try the example:  
``` 
$ cd examples 
$ numand &          # Start the server
Starting gRPC user service on [::]:50051...

$ numan             # Run the client 
General Usage:-
        command <param1> [param2] [..]. 
                Syntax: <mandatory> , [optional]

# Sample usage (see /scripts folder for more examples)
$ numan summary                 # prints a summary of the example database 
$ numan list 353-01             # prints details of numbers starting 353-01
$ numan view 353-01-12345111    # view all details for numer 353-01-12345111
```

## Usage 

```
General Usage:-
        numan command <param1> [param2] [..]. 
                Syntax: <mandatory> , [optional] 

Supported Commands:-
        summary
                Provides a summary of number database

        add <phonenumber> <domain> <carrier>
                Adds a new number to the database. Number format is cc-ndc-sn

        list_free <phonenumber> [domain] 
                Lists available numbers in db entries matching a number search. Number format is cc-ndc-sn, partial numbers are accepted

        view <phonenumber>
                Views all details and history for number entries matching a number search. Number format is cc-ndc-sn, partial numbers are accepted

        reserve <phonenumber> <oid> <minutes>
                Reserves a number for an owner for a number of minutes

        portout <phonenumber> <date>
                Sets a porting out date (dd/mm/yy)

        deallocate <phonenumber>
                De-allocates a number from an owner

        history <phonenumber>
                Lists history log for a number

        history_owner <oid>
                Lists history log for an owner

        list <phonenumber> [domain] 
                Lists number db entries matching a number search. Number format is cc-ndc-sn, partial numbers are accepted

        list_owner <oid>
                  Lists numbers attached to owner

        delete <phonenumber>
                Deletes a number permentantly (history retained)

        portin <phonenumber> <date>
                Sets a porting in date (dd/mm/yy)

        allocate <phonenumber> <oid>
                Allocates a number to an owner

```  


### Runtime Problems

#### 1. You get unusual characters in command printout (as shown below).
A terminal which supports ANSI colors is required. If the command printout looks something like below, then your terminal is not supporting ANSI colours/escape sequences correctly. 
```
←[97;40m
General Usage:-←[0m
......
```
#### 2. Message: 'Failed to load required environmental variables for config'.
You need to set environmental variables or read from a file. See examples folder. 

#### 3. Message: 'Authentication error x509: certificate signed by unknown authority"
Check the config environmental variable TLS_CERT. If using a self-signed cert then this must be set. 
See section on running client-server above for more detials.

## General Application Requirements 
- remote API & command interface
- client role based authentication
- able to reserve/hold numbers (time)
- allocate/limit numbers to a particular url domain
- select a list of random available numbers 
- able to mark numbers used. tie to account & url domain 
- have a period of quarantine when number is free  
- see number block owner/provider 
- load in batches / individually 
- remove numbers 
- log all history (number/owner/user)
- number search with wild card. 
- single owner per number 


## TODO
- command view for single number with history
- extend auth to all methods
- improve error handling 
- expand tests 
- add command user roles 
- memory store for user auth
- sanity check/verification of user/pass 
- add/remove call for users 
- add/remove users command 

## Project folder structure 
```
/   #root contains the services 'schema' (structs/interfaces). 
    /cmd
        /numand     # server 
        /numan      # command line client 
    /internal 
        /service        # core service applications 
                /datastore    # db storage layer (sqlite in this case)
        /cmdcli     # simple cli helper lib 
     /scripts       # external scripts 
    /api
        /grpc       # gRPC protobuff def & generated files 
    /examples       # example installation
#Not implemented but may be added later
    /vendor #Application dependencies (go mod controls this)
    /configs #default configs.
    /init #System init (systemd, upstart, sysv)
    /test #Additional external test services and test data.
    /docs #Design and user documents (in addition to your godoc generated documentation).
```


## API 

### internal 
Once numand is installed you can explore the api with [grpcurl]( https://github.com/fullstorydev/grpcurl
) (uses reflection)
```
# All services (-insecure required if using CA is not recognised)
grpcurl -insecure localhost:50051 describe

#or specific service 
grpcurl -insecure localhost:50051 describe grpc.Numbering

#or method
grpcurl -insecure localhost:50051 describe grpc.Numbering.Add

#or message 
grpcurl -insecure localhost:50051 describe grpc.AddRequest
```


### external 
No external API's at this time. 
public facing API could show free numbers for example. 
would require locking mechanism for reservation
would require prevention of 'mass booking', perhaps client lock. 

## Useful Links
* Sqlite command line tools - https://www.sqlite.org/cli.html
* grpcurl - https://github.com/fullstorydev/grpcurl

## Useful References

* https://grpc.io/docs/languages/go/
* https://developers.google.com/protocol-buffers
* https://github.com/golang-standards/project-layout
* https://github.com/neocortical/mysvc
* https://github.com/benbjohnson/wtf/tree/daadc79f3778fd49db6e4064878030487e2e2a47
* https://dev.to/techschoolguru/use-grpc-interceptor-for-authorization-with-jwt-1c5h
* https://medium.com/@nate510/structuring-go-grpc-microservices-dd176fdf28d0
* https://medium.com/@amsokol.com/tutorial-how-to-develop-go-grpc-microservice-with-http-rest-endpoint-middleware-kubernetes-daebb36a97e9


