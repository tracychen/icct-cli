package commands

import (
	"context"
	"fmt"
	"math/big"

	erc20tokenremote "github.com/ava-labs/avalanche-interchain-token-transfer/abi-bindings/go/TokenRemote/ERC20TokenRemote"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tracychen/ictt-cli/utils"
)

func SendTokensRemoteToHome(
	remoteRPCURL string,
	remoteEVMChainID *big.Int,
	erc20TokenAddress string,
	privateKey string,
	tokenRemoteAddress common.Address,
	amountInEther int64,
) error {
	senderKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return err
	}
	_, senderAddress, err := utils.GetPublicKeyAndAddress(senderKey)
	if err != nil {
		return err
	}
	fmt.Println("Using address as sender:", senderAddress)

	remoteOpts, err := bind.NewKeyedTransactorWithChainID(
		senderKey,
		remoteEVMChainID,
	)
	if err != nil {
		return err
	}
	remoteRc, err := rpc.Dial(remoteRPCURL)
	if err != nil {
		return err
	}
	remoteETHClient := ethclient.NewClient(remoteRc)

	tokenRemote, err := erc20tokenremote.NewERC20TokenRemote(tokenRemoteAddress, remoteETHClient)
	if err != nil {
		return err
	}
	amount := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(amountInEther))

	fmt.Println("Approving tokens on remote chain...")
	approveTx, err := tokenRemote.Approve(
		remoteOpts,
		tokenRemoteAddress,
		amount,
	)

	if err != nil {
		return err
	}

	fmt.Println("Waiting for approval transaction to be mined, tx hash:", approveTx.Hash().Hex())
	_, err = utils.WaitMinedSuccess(context.Background(), remoteETHClient, approveTx)
	if err != nil {
		return err
	}

	homeBlockchainID, err := tokenRemote.TokenHomeBlockchainID(&bind.CallOpts{})
	if err != nil {
		return err
	}
	fmt.Println("Retrieved home blockchain ID from remote:", homeBlockchainID)

	tokenHomeAddress, err := tokenRemote.TokenHomeAddress(&bind.CallOpts{})
	if err != nil {
		return err
	}

	fmt.Println("Sending", amountInEther, "tokens from remote to home")
	sendTx, err := tokenRemote.Send(
		remoteOpts,
		erc20tokenremote.SendTokensInput{
			DestinationBlockchainID:           homeBlockchainID,
			DestinationTokenTransferrerAddress: tokenHomeAddress,
			Recipient:                          common.HexToAddress(senderAddress),
			PrimaryFeeTokenAddress:             common.HexToAddress(erc20TokenAddress),
			PrimaryFee:                         big.NewInt(0),
			SecondaryFee:                       big.NewInt(0),
			RequiredGasLimit:                   big.NewInt(100_000),
		},
		amount,
	)

	if err != nil {
		return err
	}
	
	fmt.Println("Waiting for send transaction to be mined, tx hash:", sendTx.Hash().Hex())
	_, err = utils.WaitMinedSuccess(context.Background(), remoteETHClient, sendTx)
	if err != nil {
		return err
	}

	fmt.Println("Tokens sent from remote to home")
	return nil
}