package main

import (
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tracychen/ictt-cli/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
        Name:  "ictt-cli",
        Usage: "A CLI for interacting with Avalanche Interchain Token Transfer contracts",
		Commands: []*cli.Command{
			{
                Name:    "deploy-token-home",
                Aliases: []string{"dth"},
                Usage:   "Deploy a TokenHome contract",
                Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "rpc-url",
						Usage: "The RPC URL of the home chain",
						Value: "https://api.avax-test.network/ext/bc/C/rpc",
					},
					&cli.StringFlag{
						Name:  "teleporter-registry-address",
						Usage: "The address of the TeleporterRegistry contract on the home chain",
						Value: "0xF86Cb19Ad8405AEFa7d09C778215D2Cb6eBfB228",
					},
					&cli.Int64Flag{
						Name:  "chain-id",
						Usage: "The EVM chain ID of the home chain",
						Value: 43113,
					},
					&cli.StringFlag{
						Name: "token-address",
						Usage: "The address of the original ERC20 token on the home chain",
						Value: "0x7dd70607233D843b4D3EE70732f012E5985a2c78",
					},
					&cli.StringFlag{
						Name: "private-key",
						Usage: "The private key of the deployer",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					homeRPCURL := cCtx.String("rpc-url")
					homeTeleporterRegistryAddress := common.HexToAddress(cCtx.String("teleporter-registry-address"))
					homeChainID := cCtx.Int64("chain-id")
					erc20TokenAddress := common.HexToAddress(cCtx.String("token-address"))
					privateKey := cCtx.String("private-key")
                    commands.DeployTokenHome(
						homeRPCURL,
						big.NewInt(homeChainID),
						homeTeleporterRegistryAddress,
						erc20TokenAddress,
						privateKey,
					)
					return nil
                },
            },
            {
                Name:    "deploy-token-remote",
                Aliases: []string{"dtr"},
                Usage:   "Deploy a TokenRemote contract, defaults to use existing Fuji C-Chain TokenHome contract for ERC20 0x7dd70607233D843b4D3EE70732f012E5985a2c78",
                Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "home-rpc-url",
						Usage: "The RPC URL of the home chain",
						Value: "https://api.avax-test.network/ext/bc/C/rpc",
					},
					&cli.StringFlag{
						Name:  "remote-rpc-url",
						Usage: "The RPC URL of the remote chain",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "teleporter-registry-address",
						Usage: "The address of the TeleporterRegistry contract on the remote chain",
						Required: true,
					},
					&cli.Int64Flag{
						Name:  "chain-id",
						Usage: "The EVM chain ID of the remote chain",
						Required: true,
					},
					&cli.StringFlag{
						Name: "private-key",
						Usage: "The private key of the deployer",
						Required: true,
					},
					&cli.StringFlag{
						Name: "token-home",
						Usage: "The address of the TokenHome contract",
						Value: "0x8245dEb34034d98cfe941C03f1a28fE8CC778FE6",
					},
					&cli.StringFlag{
						Name: "token-address",
						Usage: "The address of the original ERC20 token",
						Value: "0x7dd70607233D843b4D3EE70732f012E5985a2c78",
					},
				},
				Action: func(cCtx *cli.Context) error {
					homeRPCURL := cCtx.String("home-rpc-url")
					remoteRPCURL := cCtx.String("remote-rpc-url")
					remoteTeleporterRegistryAddress := common.HexToAddress(cCtx.String("teleporter-registry-address"))
					remoteChainID := cCtx.Int64("chain-id")
					tokenHomeAddress := common.HexToAddress(cCtx.String("token-home"))
					erc20TokenAddress := common.HexToAddress(cCtx.String("token-address"))
					privateKey := cCtx.String("private-key")
                    commands.DeployTokenRemote(
						homeRPCURL,
						remoteRPCURL,
						remoteTeleporterRegistryAddress,
						big.NewInt(remoteChainID),
						erc20TokenAddress,
						tokenHomeAddress,
						privateKey,
					)
					return nil
                },
            },
			{
				Name:    "send-home-to-remote",
				Aliases: []string{"shr"},
				Usage:   "Send tokens from the home chain to the remote chain, defaults to Fuji C-Chain as home",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "rpc-url",
						Usage: "The RPC URL of the home chain",
						Value: "https://api.avax-test.network/ext/bc/C/rpc",
					},
					&cli.Int64Flag{
						Name:  "chain-id",
						Usage: "The EVM chain ID of the home chain",
						Value: 43113,
					},
					&cli.StringFlag{
						Name: "remote-blockchain-id",
						Usage: "The blockchain ID of the remote chain",
						Required: true,
					},
					&cli.StringFlag{
						Name: "token-address",
						Usage: "The address of the original ERC20 token",
						Value: "0x7dd70607233D843b4D3EE70732f012E5985a2c78",
					},
					&cli.StringFlag{
						Name: "private-key",
						Usage: "The private key of the sender",
						Required: true,
					},
					&cli.StringFlag{
						Name: "token-home",
						Usage: "The address of the TokenHome contract",
						Value: "0x8245dEb34034d98cfe941C03f1a28fE8CC778FE6",
					},
					&cli.StringFlag{
						Name: "token-remote",
						Usage: "The address of the TokenRemote contract",
						Required: true,
					},
					&cli.Int64Flag{
						Name: "amount",
						Usage: "The amount of tokens to send",
						Value: 1,
					},
				},
				Action: func(cCtx *cli.Context) error {
					sourceRPCURL := cCtx.String("rpc-url")
					sourceEVMChainID := cCtx.Int64("chain-id")
					remoteBlockchainID := cCtx.String("remote-blockchain-id")
					erc20TokenAddress := common.HexToAddress(cCtx.String("token-address"))
					privateKey := cCtx.String("private-key")
					tokenHomeAddress := common.HexToAddress(cCtx.String("token-home"))
					tokenRemoteAddress := common.HexToAddress(cCtx.String("token-remote"))
					amountInEther := cCtx.Int64("amount")
					commands.SendTokensHomeToRemote(
						sourceRPCURL,
						big.NewInt(sourceEVMChainID),
						remoteBlockchainID,
						erc20TokenAddress,
						privateKey,
						tokenHomeAddress,
						tokenRemoteAddress,
						amountInEther,
					)
					return nil
				},
			},
			{
				Name:    "send-remote-to-home",
				Aliases: []string{"srh"},
				Usage:   "Send tokens from the remote chain to the home chain, defaults to Fuji C-Chain as home",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "rpc-url",
						Usage: "The RPC URL of the remote chain",
						Required: true,
					},
					&cli.Int64Flag{
						Name:  "chain-id",
						Usage: "The EVM chain ID of the remote chain",
						Required: true,
					},
					&cli.StringFlag{
						Name: "token-address",
						Usage: "The address of the original ERC20 token",
						Value: "0x7dd70607233D843b4D3EE70732f012E5985a2c78",
					},
					&cli.StringFlag{
						Name: "private-key",
						Usage: "The private key of the sender",
						Required: true,
					},
					&cli.StringFlag{
						Name: "token-remote",
						Usage: "The address of the TokenRemote contract",
						Required: true,
					},
					&cli.Int64Flag{
						Name: "amount",
						Usage: "The amount of tokens to send",
						Value: 1,
					},
				},
				Action: func(cCtx *cli.Context) error {
					remoteRPCURL := cCtx.String("rpc-url")
					remoteChainID := cCtx.Int64("chain-id")
					tokenAddress := cCtx.String("token-address")
					privateKey := cCtx.String("private-key")
					tokenRemoteAddress := common.HexToAddress(cCtx.String("token-remote"))
					amountInEther := cCtx.Int64("amount")
					commands.SendTokensRemoteToHome(
						remoteRPCURL,
						big.NewInt(remoteChainID),
						tokenAddress,
						privateKey,
						tokenRemoteAddress,
						amountInEther,
					)
					return nil
				},
			},
		},
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}