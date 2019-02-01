package fetcher

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/WeTrustPlatform/account-indexer/core/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethCommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

type MockEthClient struct {
}

var header = &gethtypes.Header{
	Number: big.NewInt(2019),
}

var fromStr = "address From"
var toStr = "address To"
var toStr2 = "address To2"
var from = gethCommon.BytesToAddress([]byte(fromStr))
var to = gethCommon.BytesToAddress([]byte(toStr))
var to2 = gethCommon.BytesToAddress([]byte(toStr2))

var amount = big.NewInt(100)

var transactions = []*gethtypes.Transaction{
	gethtypes.NewTransaction(uint64(1), to, amount, uint64(21000), big.NewInt(100), nil),
	gethtypes.NewTransaction(uint64(1), to2, amount, uint64(21000), big.NewInt(100), nil),
}

func (mec MockEthClient) SubscribeNewHead(ctx context.Context, ch chan<- *gethtypes.Header) (ethereum.Subscription, error) {
	go func() {
		time.Sleep(time.Second * 2)
		ch <- header
	}()
	return nil, nil
}

func (mec MockEthClient) BlockByNumber(ctx context.Context, number *big.Int) (*gethtypes.Block, error) {
	block := gethtypes.NewBlockWithHeader(header).WithBody(transactions, nil)
	return block, nil
}

func (mec MockEthClient) TransactionSender(ctx context.Context, tx *gethtypes.Transaction, block common.Hash, index uint) (common.Address, error) {
	return from, nil
}

func (mec MockEthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*gethtypes.Header, error) {
	return header, nil
}

var count = 0

func (mec MockEthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*gethtypes.Receipt, error) {
	count++
	// 1st transaction has good tx receipt
	if count == 1 {
		return &gethtypes.Receipt{Status: 1}, nil
	}
	// 2nd transaction has bad tx receipt
	// Should not include this in the result
	return &gethtypes.Receipt{Status: 0}, nil

}

func (mec MockEthClient) TransactionByHash(ctx context.Context, hash common.Hash) (*gethtypes.Transaction, bool, error) {
	return &gethtypes.Transaction{}, false, nil
}

func (mec MockEthClient) Close() {

}

func TestFetchData(t *testing.T) {
	fetcher := ChainFetch{
		Client: MockEthClient{},
	}
	indexerChannel := make(chan *types.BLockDetail)
	go fetcher.RealtimeFetch(indexerChannel)
	// time.Sleep(time.Second * 1)
	blockDetail := <-indexerChannel
	assert.Equal(t, header.Number.Uint64(), blockDetail.BlockNumber.Uint64())
	assert.Equal(t, 1, len(blockDetail.Transactions))
	transaction := blockDetail.Transactions[0]
	assert.Equal(t, from.String(), transaction.From)
	assert.Equal(t, to.String(), transaction.To)
}
