package rpc

import (
	"fmt"
	"sync/atomic"
	"time"
)

// IvaTransactionHandler is a function type for processing transactions
type IvaTransactionHandler interface {
    ProcessIvaTransactions(transaction VMDataTransaction)
}

// IvaTransactionSubscription handles subscription to IVA transactions
type IvaTransactionSubscription struct {
    vmID           int
    startingBlock  int
    handler        IvaTransactionHandler
    pollInterval   int
    isRunning      atomic.Bool
    isPaused       atomic.Bool
    isStopped      atomic.Bool
}

// NewIvaTransactionSubscription creates a new subscription instance
func NewIvaTransactionSubscription(
    vmID int,
    startingBlock int,
    handler IvaTransactionHandler,
    pollInterval int,
) *IvaTransactionSubscription {
    return &IvaTransactionSubscription{
        vmID:          vmID,
        startingBlock: startingBlock,
        handler:       handler,
        pollInterval:  pollInterval,
    }
}

// Start begins the subscription process
func (s *IvaTransactionSubscription) Start() error {
    if s.isRunning.Load() {
        return fmt.Errorf("IvaTransactionSubscription is already running")
    }

    s.isRunning.Store(true)
    s.isPaused.Store(false)
    s.isStopped.Store(false)

    currentBlock := s.startingBlock

    for !s.isStopped.Load() {
        if s.isPaused.Load() {
            continue
        }

        latestBlock := GetLatestBlockNumber()

        effectiveLatestBlock := latestBlock
        if latestBlock > currentBlock+1000 {
            effectiveLatestBlock = currentBlock + 1000
        }

        if effectiveLatestBlock >= currentBlock {
            transactions := GetVmDataTransactions(currentBlock, effectiveLatestBlock, s.vmID)

            for _, tx := range transactions {
                s.handler.ProcessIvaTransactions(tx)
            }
            currentBlock = effectiveLatestBlock + 1
        }

        time.Sleep(time.Millisecond * 100)
    }

    s.isRunning.Store(false)
    return nil
}

func (s *IvaTransactionSubscription) Pause() {
    s.isPaused.Store(true)
}

func (s *IvaTransactionSubscription) Resume() {
    s.isPaused.Store(false)
}

func (s *IvaTransactionSubscription) Stop() {
    s.isStopped.Store(true)
}

func (s *IvaTransactionSubscription) IsRunning() bool {
    return s.isRunning.Load()
}

func (s *IvaTransactionSubscription) IsPaused() bool {
    return s.isPaused.Load()
}

func (s *IvaTransactionSubscription) IsStopped() bool {
    return s.isStopped.Load()
}

func (s *IvaTransactionSubscription) GetStartingBlock() int {
    return s.startingBlock
}

func (s *IvaTransactionSubscription) GetVmID() int {
    return s.vmID
}

func (s *IvaTransactionSubscription) GetHandler() IvaTransactionHandler {
    return s.handler
}
