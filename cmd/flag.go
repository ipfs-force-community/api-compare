package cmd

import "github.com/urfave/cli/v2"

var LotusURLFlag = &cli.StringFlag{
	Name:  "lotus-url",
	Value: "/ip4/192.168.25.124/tcp/1234",
	Usage: "lotus url",
}
var LotusTokenFlag = &cli.StringFlag{
	Name:  "lotus-token",
	Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwZXJtIjoiYWRtaW4iLCJleHQiOiIifQ.YHIIXispL9G-foNaK2PLvvqzuNPVrgRBPuLpGXymbZ0",
	Usage: "lotus token",
}

// https://api.node.glif.io
var VenusURLFlag = &cli.StringFlag{
	Name:  "venus-url",
	Value: "/ip4/192.168.25.124/tcp/3453",
	Usage: "venus url",
}
var VenusTokenFlag = &cli.StringFlag{
	Name:  "venus-token",
	Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJwZXJtIjoiYWRtaW4iLCJleHQiOiIifQ.YHIIXispL9G-foNaK2PLvvqzuNPVrgRBPuLpGXymbZ0",
	Usage: "venus token",
}

var WalletURLFlag = &cli.StringFlag{
	Name:  "wallet-url",
	Value: "/ip4/127.0.0.1/tcp/5678",
}

var WalletTokenFlag = &cli.StringFlag{
	Name:  "wallet-token",
	Value: "",
}

var GasFeeCapFlag = &cli.StringFlag{
	Name:  "gasfeecap",
	Usage: "eg. 1,1000000afil",
	Value: "1000000afil",
}
