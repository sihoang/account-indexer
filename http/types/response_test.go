package types

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/WeTrustPlatform/account-indexer/common"
	coreTypes "github.com/WeTrustPlatform/account-indexer/core/types"
	"github.com/stretchr/testify/assert"
)

var index = coreTypes.AddressIndex{
	AddressSequence: coreTypes.AddressSequence{
		Address:  "from1",
		Sequence: 1,
	},
	TxHash: "0xtx1",
	Value:  big.NewInt(-111),
	Time:   big.NewInt(1546848896),
	// BlockNumber:   big.NewInt(2018),
	CoupleAddress: "to1",
}

func TestMarshall(t *testing.T) {
	idx := AddressToEIAddress(index)
	idx.Data = []byte{1, 2}
	response := EITransactionsByAccount{
		Total:   10,
		Start:   5,
		Indexes: []EIAddress{idx},
	}
	data, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Nil(t, err)
	dataStr := string(data)
	log.Printf("%v \n", dataStr)
	tm := common.UnmarshallIntToTime(big.NewInt(1546848896)).Format(time.RFC3339)
	expectedStr := fmt.Sprintf(`{"numFound":10,"start":5,"data":[{"address":"from1","txHash":"0xtx1","value":-111,"time":"%v","coupleAddress":"to1","data":"AQI=","gas":0,"gasPrice":null}]}`, tm)
	assert.Equal(t, expectedStr, dataStr)
}

func TestUnmarshall(t *testing.T) {
	// http://mainnet.kivutar.me:3000/api/v1/accounts/0x7C419672d84a53B0a4eFed57656Ba5e4A0379084?rows=42
	str := `{"numFound":42,"start":0,"data":[{"address":"0x7c419672d84a53b0a4efed57656ba5e4a0379084","sequence":0,"tx_hash":"0x30a012177578ace1bc66cdda7d537f874682617182f41e46af439b06a5c4b9bb","value":1299152264762256509567140741597411934085535681791137795914250996166704,"time":"2019-01-07T15:14:56+07:00","coupleAddress":"0x94c9f3b353f215e8db0bef305f25e010e2441f2d","data":null,"gas":0,"gasPrice":null}]}`
	var result EITransactionsByAccount
	err := json.Unmarshal([]byte(str), &result)
	assert.Equal(t, 42, result.Total)
	assert.Nil(t, err)
}
