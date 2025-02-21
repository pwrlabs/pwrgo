package rpc

import (
	"fmt"
	"sync/atomic"
	"time"
)

// VidaTransactionHandler is a function type for processing transactions
type VidaTransactionHandler interface {
	ProcessVidaTransactions(transaction VMDataTransaction)
}

// VidaTransactionSubscription handles subscription to IVA transactions
type VidaTransactionSubscription struct {
	vmID          int
	startingBlock int
	handler       VidaTransactionHandler
	pollInterval  int
	isRunning     atomic.Bool
	isPaused      atomic.Bool
	isStopped     atomic.Bool
}

// NewVidaTransactionSubscription creates a new subscription instance
func NewVidaTransactionSubscription(
	vmID int,
	startingBlock int,
	handler VidaTransactionHandler,
	pollInterval int,
) *VidaTransactionSubscription {
	return &VidaTransactionSubscription{
		vmID:          vmID,
		startingBlock: startingBlock,
		handler:       handler,
		pollInterval:  pollInterval,
	}
}

// Start begins the subscription process
func (s *VidaTransactionSubscription) Start() error {
	if s.isRunning.Load() {
		return fmt.Errorf("VidaTransactionSubscription is already running")
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
				s.handler.ProcessVidaTransactions(tx)
			}
			currentBlock = effectiveLatestBlock + 1
		}

		time.Sleep(time.Millisecond * 100)
	}

	s.isRunning.Store(false)
	return nil
}

func (s *VidaTransactionSubscription) Pause() {
	s.isPaused.Store(true)
}

func (s *VidaTransactionSubscription) Resume() {
	s.isPaused.Store(false)
}

func (s *VidaTransactionSubscription) Stop() {
	s.isStopped.Store(true)
}

func (s *VidaTransactionSubscription) IsRunning() bool {
	return s.isRunning.Load()
}

func (s *VidaTransactionSubscription) IsPaused() bool {
	return s.isPaused.Load()
}

func (s *VidaTransactionSubscription) IsStopped() bool {
	return s.isStopped.Load()
}

func (s *VidaTransactionSubscription) GetStartingBlock() int {
	return s.startingBlock
}

func (s *VidaTransactionSubscription) GetVmID() int {
	return s.vmID
}

func (s *VidaTransactionSubscription) GetHandler() VidaTransactionHandler {
	return s.handler
}
