module github.com/footfish/numan

go 1.15

replace github.com/footfish/numan => ./

replace github.com/footfish/numan/storage => ./storage

require (
	github.com/gookit/color v1.3.7
	golang.org/x/sys v0.0.0-20210301091718-77cc2087c03b // indirect
	modernc.org/sqlite v1.8.7
)
