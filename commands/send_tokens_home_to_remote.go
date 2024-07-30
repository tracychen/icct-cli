package commands

import (
	"context"
	"fmt"
	"math/big"

	erc20tokenhome "github.com/ava-labs/avalanche-interchain-token-transfer/abi-bindings/go/TokenHome/ERC20TokenHome"
	exampleerc20 "github.com/ava-labs/avalanche-interchain-token-transfer/abi-bindings/go/mocks/ExampleERC20Decimals"
	"github.com/ava-labs/avalanchego/ids"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/ethclient"
	"github.com/ava-labs/subnet-evm/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tracychen/ictt-cli/utils"
)

func SendTokensHomeToRemote(
	homeRPCURL string,
	homeEVMChainID *big.Int,
	remoteBlockchainID string,
	erc20TokenAddress common.Address,
	privateKey string,
	tokenHomeAddress common.Address,
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

	homeOpts, err := bind.NewKeyedTransactorWithChainID(
		senderKey,
		homeEVMChainID,
	)
	if err != nil {
		return err
	}

	// Read decimals from the ERC20 token
	rc, err := rpc.Dial(homeRPCURL)
	if err != nil {
		return err
	}
	homeETHClient := ethclient.NewClient(rc)

	erc20Token, err := exampleerc20.NewExampleERC20Decimals(erc20TokenAddress, homeETHClient)
	if err != nil {
		return err
	}

	destBlockchainID, err := ids.FromString(remoteBlockchainID)
	if err != nil {
		return err
	}
	// Send tokens from C-Chain to recipient on subnet A
	input := erc20tokenhome.SendTokensInput{
		DestinationBlockchainID:           destBlockchainID,
		DestinationTokenTransferrerAddress: tokenRemoteAddress,
		Recipient:                          common.HexToAddress(senderAddress),
		PrimaryFeeTokenAddress:             erc20TokenAddress,
		PrimaryFee:                         big.NewInt(0),
		SecondaryFee:                       big.NewInt(0),
		RequiredGasLimit:                   big.NewInt(100_000),
	}

	amount := new(big.Int).Mul(big.NewInt(1e18), big.NewInt(amountInEther))

	fmt.Println("Approving tokens on home chain...")
	tx, err := erc20Token.Approve(homeOpts, tokenHomeAddress, amount)
	if err != nil {
		return err
	}

	fmt.Println("Waiting for approval transaction to be mined, tx hash:", tx.Hash().Hex())
	_, err = utils.WaitMinedSuccess(context.Background(), homeETHClient, tx)
	if err != nil {
		return err
	}

	tokenHome, err := erc20tokenhome.NewERC20TokenHome(tokenHomeAddress, homeETHClient)
	if err != nil {
		return err
	}

	fmt.Println("Sending", amountInEther, "tokens from home to remote")
	sendTx, err := tokenHome.Send(
		homeOpts,
		input,
		amount,
	)
	if err != nil {
		return err
	}

	fmt.Println("Waiting for send transaction to be mined, tx hash:", sendTx.Hash().Hex())
	_, err = utils.WaitMinedSuccess(context.Background(), homeETHClient, sendTx)
	if err != nil {
		return err
	}

	fmt.Println("Tokens sent from home to remote")
	return nil
}