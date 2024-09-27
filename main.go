package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	Token "eth-manage/token"
)

var (
	keystoreDir      string
	infuraKey        string
	network          string
	keystorePassword string
	ethNodeURL       string
	chainId          big.Int
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Read values from environment variables
	keystoreDir = os.Getenv("KESTORE_DIR")
	infuraKey = os.Getenv("INFURA_KEY")
	network = os.Getenv("NETWORK")
	keystorePassword = os.Getenv("KEYSTORE_PASSWORD")

	chainId = *big.NewInt(1)
	ethNodeURL = fmt.Sprintf("https://%s.infura.io/v3/%s", network, infuraKey)

	app := &cli.App{
		Name:  "eth_project",
		Usage: "Ethereum CLI project",
		Commands: []*cli.Command{
			{
				Name:   "create-account",
				Usage:  "Create a new Ethereum account",
				Action: createAccount,
			},
			{
				Name:   "list-accounts",
				Usage:  "List all Ethereum accounts",
				Action: listAccounts,
			},
			{
				Name:   "check-balance",
				Usage:  "Check ETH and token balances",
				Action: checkBalance,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "index",
						Usage:    "Account index to check balance",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "token-address",
						Usage:    "Token address",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "decimal",
						Usage:    "Token decimal",
						Required: false,
						Value:    6,
					},
				},
			},
			{
				Name:   "transfer-eth",
				Usage:  "Transfer ETH to another address",
				Action: transferEth,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "from",
						Usage:    "Index of the sending account",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "to",
						Usage:    "Recipient address",
						Required: true,
					},
					&cli.Float64Flag{
						Name:     "amount",
						Usage:    "Amount of ETH to transfer",
						Required: true,
					},
				},
			},
			{
				Name:   "transfer-token",
				Usage:  "Transfer tokens to another address",
				Action: transferToken,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "from",
						Usage:    "Index of the sending account",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "to",
						Usage:    "Recipient address",
						Required: true,
					},
					&cli.Float64Flag{
						Name:     "amount",
						Usage:    "Amount of tokens to transfer",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "token-address",
						Usage:    "Token address",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "decimal",
						Usage:    "Token decimal",
						Required: false,
						Value:    6,
					},
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createAccount(c *cli.Context) error {
	keyStore := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)

	account, err := keyStore.NewAccount(keystorePassword)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	log.Printf("Account created: %s", account.Address.Hex())
	return nil
}

func listAccounts(c *cli.Context) error {
	keyStore := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	accounts := keyStore.Accounts()

	if len(accounts) == 0 {
		log.Println("No accounts found.")
		return nil
	}

	for i, account := range accounts {
		fmt.Printf("Index: %d, Address: %s\n", i, account.Address.Hex())
	}
	return nil
}

// FormatBigIntToDecimal converts a big.Int amount (in Wei) to a human-readable format
// based on the provided number of decimals (e.g., 18 for Ether).
func formatBigIntToDecimal(amount *big.Int, decimals int) string {
	// Create a big float from the big.Int amount
	amountFloat := new(big.Float).SetInt(amount)

	// Create a divisor based on the token's decimals (e.g., 10^18 for Ether)
	divisor := new(big.Float).SetFloat64(float64(1))
	divisor.Mul(divisor, new(big.Float).SetFloat64(float64(1e18))) // For Ether or token with 18 decimals

	// Adjust the divisor for custom token decimals
	if decimals != 18 {
		// For token decimals other than 18
		divisor.SetFloat64(1)
		divisor = divisor.Mul(divisor, new(big.Float).SetFloat64(float64(10^decimals)))
	}

	// Divide amount by the divisor to get the human-readable amount
	humanReadable := new(big.Float).Quo(amountFloat, divisor)

	// Convert the result to a string and return it
	return humanReadable.Text('f', decimals)
}

func checkBalance(c *cli.Context) error {
	index := c.Int("index")
	tokenAddress := c.String("token-address")
	decimal := c.Int("decimal")

	keyStore := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	accounts := keyStore.Accounts()

	if index < 0 || index >= len(accounts) {
		return fmt.Errorf("invalid account index")
	}

	client, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}

	ethAddress := accounts[index].Address

	// Check ETH balance
	ethBalance, err := client.BalanceAt(context.Background(), ethAddress, nil)
	if err != nil {
		return fmt.Errorf("failed to get ETH balance: %w", err)
	}
	fmt.Printf("ETH Balance of %s: %s\n", ethAddress.Hex(), formatBigIntToDecimal(ethBalance, 18))

	// Check token balance
	tokenBalance, err := getTokenBalance(client, tokenAddress, decimal, ethAddress)
	if err != nil {
		return fmt.Errorf("failed to get token balance: %w", err)
	}
	fmt.Printf("Token Balance of %s: %s\n", ethAddress.Hex(), formatBigIntToDecimal(tokenBalance, decimal))

	return nil
}

func getTokenBalance(client *ethclient.Client, tokenAddress string, decimal int, address common.Address) (*big.Int, error) {
	tokenContract, err := Token.ERCToken(tokenAddress, decimal, client)
	if err != nil {
		return nil, err
	}

	balance, err := tokenContract.BalanceOf(address.String())
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func transferEth(c *cli.Context) error {
	fromIndex := c.Int("from")
	toAddress := c.String("to")
	amount := c.Float64("amount")

	client, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}

	keyStore := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	accounts := keyStore.Accounts()

	if fromIndex < 0 || fromIndex >= len(accounts) {
		return fmt.Errorf("invalid sender account index")
	}

	// Load the account from keystore
	account := accounts[fromIndex]

	// Unlock the account
	err = keyStore.Unlock(account, keystorePassword)
	if err != nil {
		return fmt.Errorf("failed to unlock account: %w", err)
	}

	// Create transaction
	value := big.NewInt(int64(amount * 1e18)) // Convert ETH to Wei
	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	gasLimit := uint64(21000) // Gas limit for ETH transfer
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), value, gasLimit, gasPrice, nil)

	// Sign transaction
	signedTx, err := keyStore.SignTx(account, tx, &chainId)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	return nil
}

func transferToken(c *cli.Context) error {
	fromIndex := c.Int("from")
	toAddress := c.String("to")
	amount := c.Float64("amount")
	tokenAddress := c.String("token-address")
	decimal := c.Int("decimal")

	client, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}

	keyStore := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	accounts := keyStore.Accounts()

	if fromIndex < 0 || fromIndex >= len(accounts) {
		return fmt.Errorf("invalid sender account index")
	}

	// Load the account from keystore
	account := accounts[fromIndex]

	// Unlock the account
	err = keyStore.Unlock(account, keystorePassword)
	if err != nil {
		return fmt.Errorf("failed to unlock account: %w", err)
	}

	// Create token contract instance
	tokenContract, err := Token.ERCToken(tokenAddress, decimal, client)
	if err != nil {
		return fmt.Errorf("failed to create token contract: %w", err)
	}

	// Calculate the amount in Wei
	amountInWei := new(big.Int)
	amountInWei.SetString(fmt.Sprintf("%f", amount*math.Pow10(int(decimal))), 10)

	nonce, err := client.PendingNonceAt(context.Background(), account.Address)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	gasLimit := uint64(60000) // Gas limit for token transfer
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	txData, err := tokenContract.ABI.Pack("transfer", common.HexToAddress(toAddress), amountInWei)
	if err != nil {
		return fmt.Errorf("failed to pack transfer data: %w", err)
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(tokenAddress), big.NewInt(0), gasLimit, gasPrice, txData)

	// Sign transaction
	signedTx, err := keyStore.SignTx(account, tx, &chainId)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Token transfer transaction sent: %s\n", signedTx.Hash().Hex())
	return nil
}
