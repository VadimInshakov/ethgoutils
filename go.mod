module mod

go 1.12

require (
	contract v0.0.0
	github.com/ethereum/go-ethereum v1.8.27
	golang.org/x/crypto v0.0.0-20190617133340-57b3e21c3d56
	wallet v0.0.0
)

replace wallet => ./wallet

replace contract => ./contract
