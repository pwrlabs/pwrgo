package rpc

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "sync/atomic"
    "syscall"
    "time"
)

// VidaTransactionSubscription handles subscription to VIDA transactions
type VidaTransactionSubscription struct {
    rpc                *RPC
    vidaId             int
    startingBlock      int
    latestCheckedBlock int64
    handler            ProcessVidaTransactions
    pollInterval       int
    blockSaver         BlockSaver

    running      atomic.Bool
    wantsToPause atomic.Bool
    paused       atomic.Bool
    stopped      atomic.Bool
}

// NewVidaTransactionSubscription creates a new subscription instance
func (r *RPC) NewVidaTransactionSubscription(
    vidaId int,
    startingBlock int,
    handler ProcessVidaTransactions,
    pollInterval int,
    blockSaver BlockSaver, // Can be nil
) *VidaTransactionSubscription {
    subscription := &VidaTransactionSubscription{
        rpc:                r,
        vidaId:             vidaId,
        startingBlock:      startingBlock,
        latestCheckedBlock: int64(startingBlock),
        handler:            handler,
        pollInterval:       pollInterval,
        blockSaver:         blockSaver, // nil is acceptable
    }

    // Add shutdown hook for graceful cleanup
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        log.Printf("Shutting down VidaTransactionSubscription for VIDA-ID: %d", vidaId)
        subscription.Pause()
        log.Printf("VidaTransactionSubscription for VIDA-ID: %d has been stopped.", vidaId)
    }()

    return subscription
}

// Start begins the subscription process
func (s *VidaTransactionSubscription) Start() error {
    if s.running.Load() {
        log.Printf("VidaTransactionSubscription is already running")
        return fmt.Errorf("VidaTransactionSubscription is already running")
    } else {
        s.running.Store(true)
        s.wantsToPause.Store(false)
        s.stopped.Store(false)
    }

    // Set initial block to starting block - 1 (to match Java implementation)
    atomic.StoreInt64(&s.latestCheckedBlock, int64(s.startingBlock-1))

    go func() {
        // Set goroutine name for debugging (approximation in Go)
        defer func() {
            s.running.Store(false)
        }()

        for !s.stopped.Load() {
            if s.wantsToPause.Load() {
                if !s.paused.Load() {
                    s.paused.Store(true)
                }
                time.Sleep(time.Duration(s.pollInterval) * time.Millisecond)
                continue
            }

            // If we're here and paused was true, we're resuming
            if s.paused.Load() {
                s.paused.Store(false)
            }

            func() {
                defer func() {
                    if r := recover(); r != nil {
                        log.Printf("Error in VidaTransactionSubscription: %v", r)
                    }
                }()

                latestBlock := s.rpc.GetLatestBlockNumber()
                currentLatestChecked := atomic.LoadInt64(&s.latestCheckedBlock)

                // Skip if we're already at the latest block
                if int64(latestBlock) == currentLatestChecked {
                    return
                }

                // Limit batch size to 1000 blocks (matching Java implementation)
                maxBlockToCheck := min(int64(latestBlock), currentLatestChecked+1000)

                if maxBlockToCheck > currentLatestChecked {
                    transactions := s.rpc.GetVidaDataTransactions(int(currentLatestChecked+1), int(maxBlockToCheck), s.vidaId)

                    for _, tx := range transactions {
                        func() {
                            defer func() {
                                if r := recover(); r != nil {
                                    log.Printf("Failed to process VIDA transaction: %s - %v", tx.Hash, r)
                                }
                            }()
                            s.handler(tx)
                        }()
                    }

                    // Update latest checked block
                    atomic.StoreInt64(&s.latestCheckedBlock, maxBlockToCheck)

                    // Save block state if blockSaver is provided
                    if s.blockSaver != nil {
                        func() {
                            defer func() {
                                if r := recover(); r != nil {
                                    log.Printf("Failed to save latest checked block: %d - %v", maxBlockToCheck, r)
                                }
                            }()
                            if err := s.blockSaver(int(maxBlockToCheck)); err != nil {
                                log.Printf("Failed to save latest checked block: %d - %v", maxBlockToCheck, err)
                            }
                        }()
                    }
                }
            }()

            time.Sleep(time.Duration(s.pollInterval) * time.Millisecond)
        }
    }()

    return nil
}

// min helper function for Go versions that don't have it built-in
func min(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}

func (s *VidaTransactionSubscription) Pause() {
    s.wantsToPause.Store(true)

    // Wait until the subscription is actually paused
    for !s.paused.Load() {
        time.Sleep(10 * time.Millisecond)
    }
}

func (s *VidaTransactionSubscription) Resume() {
    s.wantsToPause.Store(false)
}

func (s *VidaTransactionSubscription) Stop() {
    s.Pause() // First pause the subscription
    s.stopped.Store(true)
}

func (s *VidaTransactionSubscription) IsRunning() bool {
    return s.running.Load()
}

func (s *VidaTransactionSubscription) IsPaused() bool {
    return s.wantsToPause.Load()
}

func (s *VidaTransactionSubscription) IsStopped() bool {
    return s.stopped.Load()
}

// SetLatestCheckedBlock sets the latest checked block number
func (s *VidaTransactionSubscription) SetLatestCheckedBlock(blockNumber int) {
    atomic.StoreInt64(&s.latestCheckedBlock, int64(blockNumber))
}

func (s *VidaTransactionSubscription) GetStartingBlock() int {
    return s.startingBlock
}

func (s *VidaTransactionSubscription) GetLatestCheckedBlock() int {
    return int(atomic.LoadInt64(&s.latestCheckedBlock))
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
    options ...interface{}, // Can accept blockSaver and/or pollInterval in any order
) *VidaTransactionSubscription {
    var blockSaver BlockSaver
    interval := 100 // default poll interval

    // Parse optional parameters
    for _, option := range options {
        switch v := option.(type) {
        case BlockSaver:
            blockSaver = v
        case int:
            interval = v
        case func(int) error: // Alternative way to pass BlockSaver
            blockSaver = BlockSaver(v)
        }
    }

    subscription := r.NewVidaTransactionSubscription(
        vidaId,
        startingBlock,
        handler,
        interval,
        blockSaver,
    )

    if err := subscription.Start(); err != nil {
        log.Printf("Subscription error: %v", err)
    }

    return subscription
}
