package cmd

import (
	"fmt"

	"github.com/filecoin-project/go-jsonrpc"
	v1 "github.com/filecoin-project/venus/venus-shared/api/chain/v1"
	"github.com/filecoin-project/venus/venus-shared/api/wallet"
	"github.com/urfave/cli/v2"
)

func NewVenusFullNodeRPCFromContext(cliCtx *cli.Context) (v1.FullNode, jsonrpc.ClientCloser, error) {
	url := cliCtx.String(VenusURLFlag.Name)
	token := cliCtx.String(VenusTokenFlag.Name)
	fmt.Println("venus url:", url)
	fmt.Println("venus token:", token)

	return v1.DialFullNodeRPC(cliCtx.Context, url, token, nil)
}

func NewWalletFullRPCFromContext(cliCtx *cli.Context) (wallet.IFullAPI, jsonrpc.ClientCloser, error) {
	url := cliCtx.String(WalletURLFlag.Name)
	token := cliCtx.String(WalletTokenFlag.Name)
	fmt.Println("wallet url:", url)
	fmt.Println("wallet token:", token)

	return wallet.DialIFullAPIRPC(cliCtx.Context, url, token, nil)
}
