# numan - a phone number management tool 

This is an example Go project. It's a simple command line tool to keep track of phone number allocations. 

It's a learning project. The main purpose is to:
- explore suitable Go project layout. 
- explore using gRPC with REST endpoint. 

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
- client config file 
- setup command user roles 
- number history 
- memory store for user auth
- sanity check/verification of user/pass 
- add/remove call for users 
- add/remove users command 



## Installation
This installation uses go modules so does not need to be in your $GOPATH (of course you will need go installed to compile)

```
$ git clone https://github.com/footfish/numan
$ cd numan
$ go install ./cmd/...
 
```
## Running 

### Server 


Try example: 
```
#Use example files 
$ cd example 

$ $GOPATH/bin/numand &      
```

Alternatively:
* Install certs - see [/scripts/gen_certs.sh](./scripts/gen_certs.sh)
* Set environmental vars in file _numand.env_ ([example numand.env](./examples/numand.env))

```
$ $GOPATH/bin/numand &      
```


### Client 

`$ export RPC_ADDRESS=localhost:50051 #set RPC address to use gRPC` 

`$ $GOPATH/bin/numan          #run cli`


## General Application Requirments 
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

## Useful References

https://grpc.io/docs/languages/go/
https://developers.google.com/protocol-buffers
https://github.com/golang-standards/project-layout
https://github.com/neocortical/mysvc
https://github.com/benbjohnson/wtf/tree/daadc79f3778fd49db6e4064878030487e2e2a47
https://medium.com/@nate510/structuring-go-grpc-microservices-dd176fdf28d0
https://medium.com/@amsokol.com/tutorial-how-to-develop-go-grpc-microservice-with-http-rest-endpoint-middleware-kubernetes-daebb36a97e9


