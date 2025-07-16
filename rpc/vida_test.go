package rpc

import (
	"testing"
)

func TestSubscribeToVidaTransactions_OptionalBlockSaver(t *testing.T) {
	rpc := &RPC{RpcEndpoint: "mock://test"}
	handler := func(tx VidaDataTransaction) {}

	// Test 1: No optional parameters (should use defaults)
	subscription1 := rpc.SubscribeToVidaTransactions(1, 100, handler)
	if subscription1.blockSaver != nil {
		t.Error("BlockSaver should be nil when not provided")
	}
	if subscription1.pollInterval != 100 {
		t.Errorf("Expected default poll interval 100, got %d", subscription1.pollInterval)
	}
	subscription1.Stop()

	// Test 2: Only poll interval provided
	subscription2 := rpc.SubscribeToVidaTransactions(1, 100, handler, 250)
	if subscription2.blockSaver != nil {
		t.Error("BlockSaver should be nil when not provided")
	}
	if subscription2.pollInterval != 250 {
		t.Errorf("Expected poll interval 250, got %d", subscription2.pollInterval)
	}
	subscription2.Stop()

	// Test 3: Only blockSaver provided
	blockSaver := func(blockNumber int) error { return nil }
	subscription3 := rpc.SubscribeToVidaTransactions(1, 100, handler, blockSaver)
	if subscription3.blockSaver == nil {
		t.Error("BlockSaver should not be nil when provided")
	}
	if subscription3.pollInterval != 100 {
		t.Errorf("Expected default poll interval 100, got %d", subscription3.pollInterval)
	}
	subscription3.Stop()

	// Test 4: Both blockSaver and poll interval provided
	subscription4 := rpc.SubscribeToVidaTransactions(1, 100, handler, blockSaver, 300)
	if subscription4.blockSaver == nil {
		t.Error("BlockSaver should not be nil when provided")
	}
	if subscription4.pollInterval != 300 {
		t.Errorf("Expected poll interval 300, got %d", subscription4.pollInterval)
	}
	subscription4.Stop()

	// Test 5: Parameters in different order
	subscription5 := rpc.SubscribeToVidaTransactions(1, 100, handler, 400, blockSaver)
	if subscription5.blockSaver == nil {
		t.Error("BlockSaver should not be nil when provided")
	}
	if subscription5.pollInterval != 400 {
		t.Errorf("Expected poll interval 400, got %d", subscription5.pollInterval)
	}
	subscription5.Stop()
}

func TestNewVidaTransactionSubscription_NilBlockSaver(t *testing.T) {
	rpc := &RPC{RpcEndpoint: "mock://test"}
	handler := func(tx VidaDataTransaction) {}

	// Test that nil blockSaver works
	subscription := rpc.NewVidaTransactionSubscription(1, 100, handler, 50, nil)
	if subscription.blockSaver != nil {
		t.Error("BlockSaver should be nil when passed as nil")
	}
	if subscription.GetVidaId() != 1 {
		t.Errorf("Expected vidaId 1, got %d", subscription.GetVidaId())
	}
}
