package utils

import (
	exampleerc20 "github.com/ava-labs/avalanche-interchain-token-transfer/abi-bindings/go/mocks/ExampleERC20Decimals"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ethereum/go-ethereum/common"
)

type ERC20Details struct {
	Token *exampleerc20.ExampleERC20Decimals
	Decimals uint8
	Name string
	Symbol string
}

// Function to get the token details from an ERC20 token contract
func GetERC20Details(erc20TokenAddress common.Address, ethClient ethclient.Client) (*ERC20Details, error) {
	erc20Token, err := exampleerc20.NewExampleERC20Decimals(erc20TokenAddress, ethClient)
	if err != nil {
		return nil, err
	}
	erc20TokenDecimals, err := erc20Token.Decimals(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	erc20TokenName, err := erc20Token.Name(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	erc20TokenSymbol, err := erc20Token.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	return &ERC20Details{
		Token: erc20Token,
		Decimals: erc20TokenDecimals,
		Name: erc20TokenName,
		Symbol: erc20TokenSymbol,
	}, nil
}