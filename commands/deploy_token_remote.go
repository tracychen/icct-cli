package commands

import (
	"context"
	"fmt"
	"math/big"

	erc20tokenhome "github.com/ava-labs/avalanche-interchain-token-transfer/abi-bindings/go/TokenHome/ERC20TokenHome"
	erc20tokenremote "github.com/ava-labs/avalanche-interchain-token-transfer/abi-bindings/go/TokenRemote/ERC20TokenRemote"
	tokenremote "github.com/ava-labs/avalanche-interchain-token-transfer/abi-bindings/go/TokenRemote/TokenRemote"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tracychen/ictt-cli/utils"
)

func DeployTokenRemote(
	homeRPCURL string,
	remoteRPCURL string,
	remoteTeleporterRegistryAddress common.Address,
	remoteEVMChainID *big.Int,
	erc20TokenAddress common.Address,
	tokenHomeAddress common.Address,
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

	remoteOpts, err := bind.NewKeyedTransactorWithChainID(
		senderKey,
		remoteEVMChainID,
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

	fmt.Println("Using TokenHome deployed at:", tokenHomeAddress.Hex())

	tokenHome, err := erc20tokenhome.NewERC20TokenHome(tokenHomeAddress, homeETHClient)
	if err != nil {
		return err
	}

    fmt.Println("Deploying token remote on chain ID:", remoteEVMChainID)
	remoteRc, err := rpc.Dial(remoteRPCURL)
	if err != nil {
		return err
	}
	remoteETHClient := ethclient.NewClient(remoteRc)

	tokenHomeBlockchainID, err := tokenHome.BlockchainID(&bind.CallOpts{});
	if err != nil {
		return err
	}
	tokenRemoteSettings := erc20tokenremote.TokenRemoteSettings{
		TeleporterRegistryAddress: remoteTeleporterRegistryAddress,
		TeleporterManager: common.HexToAddress(senderAddress),
		TokenHomeBlockchainID: tokenHomeBlockchainID,
		TokenHomeAddress: tokenHomeAddress,
		TokenHomeDecimals: erc20TokenDetails.Decimals,
	}
	
	tokenRemoteAddress, tx, _, err := erc20tokenremote.DeployERC20TokenRemote(
		remoteOpts,
		remoteETHClient,
		tokenRemoteSettings,
		erc20TokenDetails.Name,
		erc20TokenDetails.Symbol,
		erc20TokenDetails.Decimals,
	)
	if err != nil {
		return err
	}

	fmt.Println("Waiting for token remote deployment transaction to be mined, tx hash:", tx.Hash().Hex())
	_, err = utils.WaitMinedSuccess(context.Background(), remoteETHClient, tx)
	if err != nil {
		return err
	}

	fmt.Printf("TokenRemote deployed at: %s\n", tokenRemoteAddress.Hex())

	tokenRemote, err := tokenremote.NewTokenRemote(tokenRemoteAddress, remoteETHClient)
	if err != nil {
		return err
	}

	fmt.Println("Registering token remote with token home...")
	sendRegisterTx, err := tokenRemote.RegisterWithHome(remoteOpts, tokenremote.TeleporterFeeInfo{
		FeeTokenAddress: erc20TokenAddress,
		Amount:          big.NewInt(0),
	})
	if err != nil {
		return err
	}

	fmt.Println("Waiting for token registration transaction to be mined, tx hash:", sendRegisterTx.Hash().Hex())
	_, err = utils.WaitMinedSuccess(context.Background(), remoteETHClient, sendRegisterTx)
	if err != nil {
		return err
	}

	fmt.Println("Token remote deployed and registered with token home, token remote address:", tokenRemoteAddress.Hex())
	return nil
}