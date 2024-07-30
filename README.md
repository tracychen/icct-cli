# ICTT CLI

A CLI tool for interacting with [Avalanche Interchain Token Transfer (ICTT)](https://github.com/ava-labs/avalanche-interchain-token-transfer) contracts. This tool provides commands to deploy and interact with TokenHome and TokenRemote contracts, as well as to send tokens between chains.

## Installation

1. Clone the repository

   ```bash
    git clone https://github.com/tracychen/ictt-cli.git
    cd ictt-cli

   ```

2. Build the CLI

   ```bash
   go install
   ```

## Usage

### Deploy TokenHome Contract

Deploys a TokenHome contract on the home chain.

```bash
ictt-cli deploy-token-home [flags]
ictt-cli dth [flags]

```

Flags:

- `--rpc-url`: The RPC URL of the home chain (default: `https://api.avax-test.network/ext/bc/C/rpc`)
- `--teleporter-registry-address`: The address of the TeleporterRegistry contract on the home chain (default: `0xF86Cb19Ad8405AEFa7d09C778215D2Cb6eBfB228`)
- `--chain-id`: The EVM chain ID of the home chain (default: `43113`)
- `--token-address`: The address of the ERC20 token on the home chain (default: `0x7dd70607233D843b4D3EE70732f012E5985a2c78`)
- `--private-key`: The private key of the deployer (required)

Example:

```bash
ictt-cli deploy-token-home \
--rpc-url="https://api.avax-test.network/ext/bc/C/rpc" \
--chain-id=43113 \
--teleporter-registry-address="0xF86Cb19Ad8405AEFa7d09C778215D2Cb6eBfB228" \
--token-address="0x7dd70607233D843b4D3EE70732f012E5985a2c78" \
--private-key="ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

```

### Deploy TokenRemote Contract

Deploys a TokenRemote contract, defaults to using the existing TokenHome contract for EXMP ERC20 token on Fuji C-Chain.

```bash
ictt-cli deploy-token-remote [flags]
ictt-cli dtr [flags]
```

Flags:

- `--home-rpc-url`: The RPC URL of the home chain (default: `https://api.avax-test.network/ext/bc/C/rpc`)
- `--remote-rpc-url`: The RPC URL of the remote chain (required)
- `--teleporter-registry-address`: The address of the TeleporterRegistry contract on the remote chain (required)
- `--chain-id`: The EVM chain ID of the remote chain (required)
- `--private-key`: The private key of the deployer (required)
- `--token-home`: The address of the TokenHome contract (default: `0x8245dEb34034d98cfe941C03f1a28fE8CC778FE6`)
- `--token-address`: The address of the ERC20 token (default: `0x7dd70607233D843b4D3EE70732f012E5985a2c78`)

Example:

```bash
ictt-cli dtr \
--remote-rpc-url="https://subnets.avax.network/dexalot/testnet/rpc" \
--teleporter-registry-address="0x7be74F0b2e89b578D003ceF8cfe00dfA1A9fD705" \
--chain-id=432201 \
--private-key="ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" \
--amount=1
```

### Send Tokens From Home to Remote

Sends tokens from the home chain to the remote chain.

```bash
ictt-cli send-home-to-remote [flags]
ictt-cli shr [flags]
```

Flags:

- `--rpc-url`: The RPC URL of the home chain (default: `https://api.avax-test.network/ext/bc/C/rpc`)
- `--chain-id`: The EVM chain ID of the home chain (default: `43113`)
- `--remote-blockchain-id`: The blockchain ID of the remote chain (required)
- `--token-address`: The address of the ERC20 token (default: `0x7dd70607233D843b4D3EE70732f012E5985a2c78`)
- `--private-key`: The private key of the sender (required)
- `--token-home`: The address of the TokenHome contract (default: `0x8245dEb34034d98cfe941C03f1a28fE8CC778FE6`)
- `--token-remote`: The address of the TokenRemote contract (required)
- `--amount`: The amount of tokens to send (default: 1)

Example:

```bash
ictt-cli send-home-to-remote \
--rpc-url="https://api.avax-test.network/ext/bc/C/rpc" \
--chain-id=43113 \
--remote-blockchain-id=XuEPnCE59rtutASDPCDeYw8geQaGWwteWjkDXYLWvssfuirde \
--token-address="0x7dd70607233D843b4D3EE70732f012E5985a2c78" \
--private-key="ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" \
--token-home="0x8245dEb34034d98cfe941C03f1a28fE8CC778FE6" \
--token-remote="0x7be74F0b2e89b578D003ceF8cfe00dfA1A9fD705" \
--amount=1
```

### Send Tokens From Remote to Home

Sends tokens from the remote chain to the home chain.

```bash
ictt-cli send-remote-to-home [flags]
ictt-cli srh [flags]
```

Flags:

- `--rpc-url`: The RPC URL of the remote chain (required)
- `--chain-id`: The EVM chain ID of the remote chain (required)
- `--token-address`: The address of the ERC20 token (default: `0x7dd70607233D843b4D3EE70732f012E5985a2c78`)
- `--private-key`: The private key of the sender (required)
- `--token-remote`: The address of the TokenRemote contract (required)
- `--amount`: The amount of tokens to send (default: 1)

Example:

```bash
ictt-cli send-remote-to-home \
--rpc-url="https://subnets.avax.network/dexalot/testnet/rpc" \
--chain-id=432201 \
--token-address="0x7dd70607233D843b4D3EE70732f012E5985a2c78" \
--private-key="ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" \
--token-remote="0x7be74F0b2e89b578D003ceF8cfe00dfA1A9fD705" \
--amount=1
```
