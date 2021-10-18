# numan - a phone number management tool 

This is an example Go project. It's a simple command line tool and gRPC API to keep track of phone number allocations/reservations stored in an sqlite db. 

It's a personal learning project (to replace an excel file). The main purpose is to:
- explore suitable Go project layout. 
- explore using gRPC (and perhaps extending to expose gRPC via REST ). 

The project loosely uses domain driven design and the  [standard Go project layout](https://github.com/golang-standards/project-layout). The main business objects being a 'numbering' record and a change 'history'. 
The business objects are defined in root, then object logic is 'layered' using the root interface (I guess you could call this Go's version of method overloading ). 

To demonstrate using an example from the project folders. 

- 'numbering' is a business object 
- /numbering.go <- this defines the business object API ie. interface and structs  
- /internal/app/numbering.go <- this implements the number interface to the overall application 
- /internal/datastore/numbering.go <- this implements the number interface to storage 

This pattern allows for: 
- Easily changing/adding new layers. For example swapping to a different storage mechanism. 
- Avoids circular dependencies (which go does not allow). 

## TODO
- improve error handling 
- expand tests 
- setup command user roles 
- number history 
- memory store for user auth
- sanity check/verification of user/pass 
- add/remove call for users 
- add/remove users command 

## Installation
This installation uses Go modules so should not installed in your $GOPATH (of course you will need Go installed to compile)

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
$ numan view 353-01-12345111    # view all details for numer 353-01-12345111

```
### Running Client-Server Mode

To run in client-server mode you will need to use certificates. 
There are two approaches you can use. 
 1) Use a self-signed cert. This approach requires the cert to be loaded in the client (or switch off verification). 
 2) Use your own trusted CA with minica (preferred). This approach does NOT require the cert to be loaded in the client but it a little trickier to set up. 

 The example uses a self-signed cert. For help installing your own certs - see [/scripts/gen_certs.sh](./scripts/gen_certs.sh)

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



### Problems

#### 1. You get unusual characters in command printout (as shown below).
A terminal which supports ANSI colors is required. If the command printout looks something like below, then your terminal is not supporting ANSI colours/escape sequences correctly. 
```
←[97;40m
General Usage:-←[0m
......
```
#### 2. MessageI 'Failed to load required environmental variables for config'.
You need to set environmental variables or read from a file. See examples folder. 

#### 3. Message: 'Authentication error x509: certificate signed by unknown authority"
Check the config environmental variable TLS_CERT. If using a self-signed cert then this must be set. 
See section on running client-server above for more detials.



## General Application Requirements 
- remote API & command interface
- role based authentication
- able to reserve/hold numbers (time)
- limit numbers to a particular url domain
- select a list of random available numbers 
- able to mark numbers used. tie to account & url domain 
- have a period of quarantine when number is free  
- see number block owner/provider 
- load in batches / individually 
- remove numbers 
- log number history 
- log user history (who/what cancelled & when) 
- number search with wild card. 
- single user per number 
- log of porting 


## Project folder structure 
```
/   #root contains business domain 'schema' (structs/interfaces). 
    /cmd
        /numand     # server 
        /numan      # command line client 
    /internal 
        /app #core application 
        /cmdcli     # simple cli helper lib 
        /datastore    # db storage (sqlite in this case)
    /scripts        # external scripts 

#Not implemented but may be added 
    /api
        /grpc       # gRPC protobuff def & generated files 
    /vendor #Application dependencies (go mod controls this)
    /configs #default configs.
    /init #System init (systemd, upstart, sysv)
    /test #Additional external test apps and test data.
    /docs #Design and user documents (in addition to your godoc generated documentation).
```


## API 

### internal 
TODO

### external 
No external API's at this time. 
public facing API could show free numbers for example. 
would require locking mechanism for reservation
would require prevention of 'mass booking', perhaps client lock. 

## Useful Links
Sqlite command line tools - https://www.sqlite.org/cli.html

## Useful References

https://grpc.io/docs/languages/go/
https://developers.google.com/protocol-buffers
https://github.com/golang-standards/project-layout
https://github.com/neocortical/mysvc
https://github.com/benbjohnson/wtf/tree/daadc79f3778fd49db6e4064878030487e2e2a47
https://medium.com/@nate510/structuring-go-grpc-microservices-dd176fdf28d0
https://medium.com/@amsokol.com/tutorial-how-to-develop-go-grpc-microservice-with-http-rest-endpoint-middleware-kubernetes-daebb36a97e9


