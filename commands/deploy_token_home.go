package commands

import (
	"context"
	"fmt"
	"math/big"

	erc20tokenhome "github.com/ava-labs/avalanche-interchain-token-transfer/abi-bindings/go/TokenHome/ERC20TokenHome"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tracychen/ictt-cli/utils"
)

func DeployTokenHome(
	homeRPCURL string,
	homeEVMChainID *big.Int,
	homeTeleporterRegistryAddress common.Address,
	erc20TokenAddress common.Address,
	privateKey string,
) error {
	senderKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}
	_, senderAddress, err := utils.GetPublicKeyAndAddress(senderKey)
	if err != nil {
		return err
	}
	fmt.Println("Using address as deployer:", senderAddress)

	homeOpts, err := bind.NewKeyedTransactorWithChainID(
		senderKey,
		homeEVMChainID,
	)
	if err != nil {
		return err
	}

	rc, err := rpc.Dial(homeRPCURL)
	if err != nil {
		return err
	}
	homeETHClient := ethclient.NewClient(rc)

	erc20TokenDetails, err := utils.GetERC20Details(erc20TokenAddress, homeETHClient)
	if err != nil {
		return err
	}
	fmt.Printf("Using ERC20 token deployed at: %s\n", erc20TokenAddress)
	fmt.Printf("ERC20 token name: %s\n", erc20TokenDetails.Name)
	fmt.Printf("ERC20 token symbol: %s\n", erc20TokenDetails.Symbol)
	fmt.Printf("ERC20 token decimals: %d\n", erc20TokenDetails.Decimals)
	fmt.Println("-----")

	fmt.Println("Deploying token home...")
	tokenHomeAddress, tx, _, err := erc20tokenhome.DeployERC20TokenHome(
		homeOpts,
		homeETHClient,
		homeTeleporterRegistryAddress,
		common.HexToAddress(senderAddress),
		erc20TokenAddress,
		erc20TokenDetails.Decimals,
	)
	if err != nil {
		return err
	}

	fmt.Println("Waiting for token home deployment transaction to be mined, tx hash:", tx.Hash().Hex())
	_, err = utils.WaitMinedSuccess(context.Background(), homeETHClient, tx)
	if err != nil {
		return err
	}
	
	fmt.Printf("TokenHome deployed at: %s\n", tokenHomeAddress.Hex())
	return err
}