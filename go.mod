module mod

go 1.12

require (
	github.com/ethereum/go-ethereum v1.8.27
	github.com/rs/cors v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20190530122614-20be4c3c3ed5
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	wallet v0.0.0
)

replace wallet => ./wallet
