module github.com/footfish/numan

go 1.15

replace github.com/footfish/numan => ./

replace github.com/footfish/numan/datastore => ./datastore

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/gookit/color v1.4.2
	github.com/joho/godotenv v1.3.0
	github.com/kataras/tablewriter v0.0.0-20180708051242-e063d29b7c23 // indirect
	github.com/lensesio/tableprinter v0.0.0-20201125135848-89e81fc956e7
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/vrischmann/envconfig v1.3.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210302154924-ca353664deba // indirect
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.27.1
	modernc.org/sqlite v1.8.7
)
