module github.com/footfish/numan

go 1.15

replace github.com/footfish/numan => ./

replace github.com/footfish/numan/storage => ./storage

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/gookit/color v1.3.7
	github.com/kr/pretty v0.1.0 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210301091718-77cc2087c03b // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210302154924-ca353664deba // indirect
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	modernc.org/sqlite v1.8.7
)
