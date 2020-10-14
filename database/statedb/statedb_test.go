package statedb

import (
	"github.com/jhdriver/UWORLD/common/utils"
	"github.com/jhdriver/UWORLD/core/types"
	"os"
	"strconv"
	"testing"
)

var dbPath = "state_test.db"

func BenchmarkStateStorage_SetAccountState(b *testing.B) {
	if utils.IsExist(dbPath) {
		if err := os.Remove(dbPath); err != nil {
			b.Errorf("clear state_test.db %s", err.Error())
		}
	}

	state := NewStateStorage(dbPath)
	err := state.Open()
	if err != nil {
		b.Errorf("open state_test.db %s", err.Error())
	}

	b.N = 1000000
	for i := 0; i < b.N; i++ {
		account := &types.Account{
			Address:         []byte(strconv.Itoa(i)),
			Balance:         0,
			LockedIn:        0,
			LockedOut:       0,
			Nonce:           0,
			Time:            0,
			ConfirmedHeight: 0,
			ConfirmedNonce:  0,
			ConfirmedTime:   0,
			JournalIn:       nil,
			JournalOut:      nil,
		}
		state.SetAccountState(account)
	}
}
