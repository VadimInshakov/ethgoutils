module mod

go 1.12

require (
	contract v0.0.0
	github.com/ethereum/go-ethereum v1.8.27
	golang.org/x/crypto v0.0.0-20190617133340-57b3e21c3d56
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	mycontract v0.0.0
	wallet v0.0.0
)

replace wallet => ./wallet

replace contract => ./contract

replace mycontract => ./build
