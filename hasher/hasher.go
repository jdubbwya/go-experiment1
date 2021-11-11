package hasher

import (
	"crypto/sha512"
	"encoding/base64"
	"io"
	"sync"
	"sync/atomic"
	"time"
)


const TransactionUnknown = 0
const TransactionInProgress = 1
const TransactionComplete = 2

var TransactionReceiptPause = 5 * time.Second // 5 seconds

// Tracks the previous allocated transaction ids
var upperTransactionId int64 = 0

// stores the values of all hashed password with a concurrency safe datastore
var hashedPasswords sync.Map

func Enqueue(rawPassword string) (int64, error) {
	id := atomic.AddInt64(&upperTransactionId, 1)
	go func() {
		time.Sleep(TransactionReceiptPause)
		h512 := sha512.New()
		io.WriteString(h512, rawPassword)
		hashedPasswords.Store( id, base64.StdEncoding.EncodeToString(h512.Sum(nil)) )
	}()
	return id, nil
}


func TransactionResult( transactionId int64 ) (*string, int) {
	var entry, ok = hashedPasswords.Load(transactionId)
	if ok {
		var hashedPassword = entry.(string)
		return &hashedPassword, TransactionComplete
	}

	if transactionId < upperTransactionId {
		return nil, TransactionUnknown
	}

	return nil, TransactionInProgress
}
