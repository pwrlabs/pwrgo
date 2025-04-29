package rpc

import (
	"fmt"
	"sync/atomic"
	"time"
)

// VidaTransactionSubscription handles subscription to VIDA transactions
type VidaTransactionSubscription struct {
	rpc                *RPC
	vidaId             int
	startingBlock      int
	latestCheckedBlock int
	handler            ProcessVidaTransactions
	pollInterval       int

	running atomic.Bool
	paused  atomic.Bool
	stopped atomic.Bool
}

// NewVidaTransactionSubscription creates a new subscription instance
func (r *RPC) NewVidaTransactionSubscription(
	vidaId int,
	startingBlock int,
	handler ProcessVidaTransactions,
	pollInterval int,
) *VidaTransactionSubscription {
	return &VidaTransactionSubscription{
		rpc:                r,
		vidaId:             vidaId,
		startingBlock:      startingBlock,
		latestCheckedBlock: startingBlock,
		handler:            handler,
		pollInterval:       pollInterval,
	}
}

// Start begins the subscription process
func (s *VidaTransactionSubscription) Start() error {
	if s.running.Load() {
		fmt.Println("VidaTransactionSubscription is already running")
		return fmt.Errorf("VidaTransactionSubscription is already running")
	} else {
		s.running.Store(true)
		s.paused.Store(false)
		s.stopped.Store(false)
	}

	currentBlock := s.startingBlock

	go func() {
		for !s.stopped.Load() {
			if s.paused.Load() {
				time.Sleep(time.Duration(s.pollInterval) * time.Millisecond)
				continue
			}

			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Error in subscription: %v\n", r)
					}
				}()

				latestBlock := s.rpc.GetLatestBlockNumber()

				effectiveLatestBlock := latestBlock
				if latestBlock > currentBlock+1000 {
					effectiveLatestBlock = currentBlock + 1000
				}

				if effectiveLatestBlock >= currentBlock {
					transactions := s.rpc.GetVidaDataTransactions(currentBlock, effectiveLatestBlock, s.vidaId)

					for _, tx := range transactions {
						s.handler(tx)
					}

					s.latestCheckedBlock = effectiveLatestBlock
					currentBlock = effectiveLatestBlock + 1
				}
			}()

			time.Sleep(time.Duration(s.pollInterval) * time.Millisecond)
		}

		s.running.Store(false)
	}()

	return nil
}

func (s *VidaTransactionSubscription) Pause() {
	s.paused.Store(true)
}

func (s *VidaTransactionSubscription) Resume() {
	s.paused.Store(false)
}

func (s *VidaTransactionSubscription) Stop() {
	s.stopped.Store(true)
}

func (s *VidaTransactionSubscription) IsRunning() bool {
	return s.running.Load()
}

func (s *VidaTransactionSubscription) IsPaused() bool {
	return s.paused.Load()
}

func (s *VidaTransactionSubscription) IsStopped() bool {
	return s.stopped.Load()
}

func (s *VidaTransactionSubscription) GetStartingBlock() int {
	return s.startingBlock
}

func (s *VidaTransactionSubscription) GetLatestCheckedBlock() int {
	return s.latestCheckedBlock
}

func (s *VidaTransactionSubscription) GetVidaId() int {
	return s.vidaId
}

func (s *VidaTransactionSubscription) GetHandler() ProcessVidaTransactions {
	return s.handler
}

// SubscribeToVidaTransactions creates and starts a subscription
func (r *RPC) SubscribeToVidaTransactions(
	vidaId int,
	startingBlock int,
	handler ProcessVidaTransactions,
	pollInterval ...int,
) *VidaTransactionSubscription {
	interval := 100
	if len(pollInterval) > 0 {
		interval = pollInterval[0]
	}

	subscription := r.NewVidaTransactionSubscription(
		vidaId,
		startingBlock,
		handler,
		interval,
	)

	if err := subscription.Start(); err != nil {
		fmt.Printf("Subscription error: %v\n", err)
	}

	return subscription
}
