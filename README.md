# Ethereum CLI Project

This project is a command-line interface (CLI) for managing Ethereum accounts, checking balances, and transferring ETH and ERC20 tokens. It uses Go and the `go-ethereum` package to interact with the Ethereum blockchain.

## Features

- **Create and manage multiple Ethereum accounts**
- **Check ETH and ERC20 token balances for specific accounts**
- **Transfer ETH and ERC20 tokens between accounts**
- **Configuration through an environment file (.env)**

## Prerequisites

- **Go** (version 1.16 or higher)
- **Access to an Ethereum node** (e.g., Infura, Alchemy)
- **Basic knowledge of Ethereum and blockchain concepts**

## Setup

1. **Clone the repository**:

   ```bash
   git clone https://github.com/ls-taylor-lee/go-manage.git
   cd eth-cli-project
   ```

2. **Install dependencies**:

   ```bash
   go get github.com/joho/godotenv
   go get github.com/ethereum/go-ethereum
   go get github.com/urfave/cli/v2
   ```

3. **Create an `.env` file**:

   Create a file named `.env` in the root of the project directory with the following content:

   ```plaintext
   KESTORE_DIR=./keystore
   TOKEN_ADDRESS=0xYourTokenContract
   KEYSTORE_PASSWORD=your_secure_password
   INFURA_KEY=your_infura_key
   NETWORK=mainnet
   CHAIN_ID=1
   ```

   Replace the placeholder values with your actual configuration. Ensure that `KESTORE_DIR` points to a directory where you want to store your keystore files.

## Usage

After setting up the project and environment, you can use the CLI commands:

### Create an Ethereum Account

```bash
go run main.go create-account
```

### List All Accounts

```bash
go run main.go list-accounts
```

### Check Balances

Check the ETH and token balances for a specific account by index:

```bash
go run main.go check-balance --index 0
```

### Transfer ETH

Transfer ETH from one account to another:

```bash
go run main.go transfer-eth --from 0 --to 0xRecipientAddress --amount 0.1
```

### Transfer Tokens

Transfer ERC20 tokens from one account to another:

```bash
go run main.go transfer-token --from 0 --to 0xRecipientAddress --amount 1
```

## Notes

- Make sure to handle your keystore securely and use strong passwords.
- Gas prices and limits may need adjustment based on current network conditions.
- Ensure that you have the necessary permissions to access the Ethereum node you are using.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgements

- [go-ethereum](https://github.com/ethereum/go-ethereum): Go implementation of the Ethereum protocol.
- [urfave/cli](https://github.com/urfave/cli): A simple, fast, and fun package for building command line apps in Go.
- [joho/godotenv](https://github.com/joho/godotenv): A Go port of the Ruby dotenv library to load environment variables from `.env` files.
