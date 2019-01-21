// +build int

package int

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/WeTrustPlatform/account-indexer/common"
	"github.com/WeTrustPlatform/account-indexer/fetcher"
	"github.com/WeTrustPlatform/account-indexer/http/types"
	"github.com/WeTrustPlatform/account-indexer/indexer"
	"github.com/WeTrustPlatform/account-indexer/repository/keyvalue"
	"github.com/WeTrustPlatform/account-indexer/repository/keyvalue/dao"
	"github.com/WeTrustPlatform/account-indexer/service"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/memdb"
)

func TestContractCreation(t *testing.T) {
	// Setup
	// ipcs := []string{"wss://mainnet.kivutar.me:8546/2KT179di"}
	ipcs := []string{"wss://mainnet.infura.io/_ws"}
	t.Logf("TestContractCreation ipcs=%v \n ", ipcs)
	service.GetIpcManager().SetIPC(ipcs)
	fetcher, err := fetcher.NewChainFetch()
	assert.Nil(t, err)
	blockNumber := big.NewInt(6808718)
	// Run Test
	blockDetail, err := fetcher.FetchABlock(blockNumber)
	assert.Nil(t, err)
	// log.Println(blockDetail)
	idx := NewTestIndexer()
	isBatch := true
	idx.ProcessBlock(blockDetail, isBatch)
	// Confirm contract created tx
	contract := "0x4a6ead96974679957a17d2f9c7835a3da7ddf91d"
	fromTime, _ := common.StrToTime("2018-12-01T00:00:00")
	toTime, _ := common.StrToTime("2018-12-01T23:59:59")
	total, addressIndexes := idx.IndexRepo.GetTransactionByAddress(contract, 10, 0, fromTime, toTime)
	assert.Equal(t, 1, total)
	tx := addressIndexes[0].TxHash
	assert.True(t, strings.EqualFold("0x61278dd960415eadf11cfe17a6c38397af658e77bbdd367db70e19ee3a193bdd", tx))
	tm := common.UnmarshallIntToTime(addressIndexes[0].Time)
	t.Logf("TestContractCreation found contract creation transaction at %v \n", tm)
	// another transaction in that block
	address := "0xec3ecca662f089a1bb83c681339729fef66e22b1"
	total, addressIndexes = idx.IndexRepo.GetTransactionByAddress(address, 10, 0, fromTime, toTime)
	assert.Equal(t, 1, total)
	tx = addressIndexes[0].TxHash
	assert.True(t, strings.EqualFold("0xdd230c27a118d2707174dca6c17ca94da6f1ed1fc1e3374373c562168c1dcd67", tx))
	val := addressIndexes[0].Value
	assert.Equal(t, big.NewInt(1201532900000000000), val)
	eiAddress := types.AddressToEIAddress(addressIndexes[0])
	ba, _ := json.Marshal(eiAddress)
	t.Logf("The returned json is %v \n", string(ba))
}

func NewTestIndexer() indexer.Indexer {
	addressDB := memdb.New(comparer.DefaultComparer, 0)
	addressDAO := dao.NewMemDbDAO(addressDB)
	blockDB := memdb.New(comparer.DefaultComparer, 0)
	blockDAO := dao.NewMemDbDAO(blockDB)
	batchDB := memdb.New(comparer.DefaultComparer, 0)
	batchDAO := dao.NewMemDbDAO(batchDB)
	indexRepo := keyvalue.NewKVIndexRepo(addressDAO, blockDAO)
	batchRepo := keyvalue.NewKVBatchRepo(batchDAO)
	idx := indexer.NewIndexer(indexRepo, batchRepo, nil)
	return idx
}
