package kernel

import (
	"testing"
)

func TestNotificationCallbacks(t *testing.T) {
	var blockTipCalled bool
	var headerTipCalled bool
	var lastBlockHeight int64
	var lastHeaderHeight int64

	callbacks := &NotificationCallbacks{
		OnBlockTip: func(state SynchronizationState, index *BlockIndex, _ float64) {
			blockTipCalled = true
			lastBlockHeight = int64(index.Height())
		},
		OnHeaderTip: func(state SynchronizationState, height int64, timestamp int64, presync bool) {
			headerTipCalled = true
			lastHeaderHeight = height
		},
	}
	suite := ChainstateManagerTestSuite{
		MaxBlockHeightToImport: 5,
		NotificationCallbacks:  callbacks,
	}
	suite.Setup(t)

	if !blockTipCalled {
		t.Error("OnBlockTip callback was not called")
	}
	if lastBlockHeight != 5 {
		t.Errorf("Expected last block height 5, got %d", lastBlockHeight)
	}

	if !headerTipCalled {
		t.Error("OnHeaderTip callback was not called")
	}
	if lastHeaderHeight != 5 {
		t.Errorf("Expected last header height 5, got %d", lastHeaderHeight)
	}
}
