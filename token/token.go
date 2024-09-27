package token

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Token struct
type Token struct {
	address common.Address
	decimal int
	ABI     abi.ABI
	client  *ethclient.Client
}

// NewToken function
func ERCToken(address string, decimal int, client *ethclient.Client) (*Token, error) {
	parsedABI, err := loadTokenABI("./token/abi.json") // Update the path to abi.json
	if err != nil {
		return nil, fmt.Errorf("failed to load token ABI: %w", err)
	}

	return &Token{
		address: common.HexToAddress(address),
		decimal: decimal,
		ABI:     parsedABI,
		client:  client,
	}, nil
}

// loadTokenABI function
func loadTokenABI(filename string) (abi.ABI, error) {
	file, err := os.Open(filename)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to open ABI file: %w", err)
	}
	defer file.Close()

	abiBytes, err := io.ReadAll(file)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to read ABI file: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return abi.ABI{}, fmt.Errorf("failed to parse ABI: %w", err)
	}

	return parsedABI, nil
}

// BalanceOf method using ethclient
func (t *Token) BalanceOf(address string) (*big.Int, error) {
	addressToCheck := common.HexToAddress(address)

	// Prepare the data for the call
	data, err := t.ABI.Pack("balanceOf", addressToCheck)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data for balanceOf: %w", err)
	}

	// Call the balanceOf function from the contract
	msg := ethereum.CallMsg{
		To:   &t.address,
		Data: data,
	}

	// Use CallContract method to execute the contract method
	result, err := t.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call balanceOf: %w", err)
	}

	// Unpack the result into a *big.Int
	balance := new(big.Int).SetBytes(result)

	return balance, nil
}
