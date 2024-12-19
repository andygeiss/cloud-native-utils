package consistency_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/andygeiss/cloud-native/utils/consistency"
)

type Transaction struct {
	ID              uint64    `json:"id"`
	DebitAccountID  string    `json:"debit_account_id"`
	CreditAccountID string    `json:"credit_account_id"`
	Amount          float64   `json:"amount"`
	Timestamp       time.Time `json:"timestamp"`
}

func Benchmark_GobFileLogger(b *testing.B) {
	jsonFile := "benchmark_gob_file_logger.json"
	defer os.Remove(jsonFile)
	n := 1000000
	// Create N transactions.
	txs := make([]Transaction, n)
	for i := 0; i < n; i++ {
		txs[i] = Transaction{ID: uint64(rand.Int63())}
	}
	// Create a GOB file logger.
	logger := consistency.NewGobFileLogger[uint64, Transaction](jsonFile)
	defer logger.Close()
	b.Run(fmt.Sprintf("write %d put events into the log file", n), func(b *testing.B) {
		for i := 0; i < n; i++ {
			tx := txs[i]
			logger.WritePut(tx.ID, tx)
		}
	})
	b.Run(fmt.Sprintf("write %d delete events into the log file", n), func(b *testing.B) {
		for i := 0; i < n; i++ {
			tx := txs[i]
			logger.WriteDelete(tx.ID)
		}
	})
}

func Benchmark_JsonFileLogger(b *testing.B) {
	jsonFile := "benchmark_json_file_logger.json"
	defer os.Remove(jsonFile)
	n := 1000000
	// Create N transactions.
	txs := make([]Transaction, n)
	for i := 0; i < n; i++ {
		txs[i] = Transaction{ID: uint64(rand.Int63())}
	}
	// Create a JSON file logger.
	logger := consistency.NewJsonFileLogger[uint64, Transaction](jsonFile)
	defer logger.Close()
	b.Run(fmt.Sprintf("write %d put events into the log file", n), func(b *testing.B) {
		for i := 0; i < n; i++ {
			tx := txs[i]
			logger.WritePut(tx.ID, tx)
		}
	})
	b.Run(fmt.Sprintf("write %d delete events into the log file", n), func(b *testing.B) {
		for i := 0; i < n; i++ {
			tx := txs[i]
			logger.WriteDelete(tx.ID)
		}
	})
}
