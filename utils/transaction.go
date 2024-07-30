package utils

import (
	"context"
	"fmt"

	"github.com/ava-labs/subnet-evm/accounts/abi/bind"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/ethclient"
)

// Function to wait for a transaction to be successfully mined
func WaitMinedSuccess(ctx context.Context, client ethclient.Client, tx *types.Transaction) (*types.Receipt, error) {
	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, err
	}

	// Check if the transaction was successful
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, fmt.Errorf("transaction failed, status: %s", receipt.Status)
	}

	return receipt, nil
}