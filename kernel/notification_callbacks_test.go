package kernel

import (
	"testing"
)

type notificationTestCallbacks struct {
	lastBlockHeight  int64
	lastHeaderHeight int64
	blockTipCalled   bool
	headerTipCalled  bool
}

func (c *notificationTestCallbacks) BlockTip(state SynchronizationState, entry *BlockTreeEntry, progress float64) {
	c.blockTipCalled = true
	c.lastBlockHeight = int64(entry.Height())
}

func (c *notificationTestCallbacks) HeaderTip(state SynchronizationState, height int64, timestamp int64, presync bool) {
	c.headerTipCalled = true
	c.lastHeaderHeight = height
}

func (c *notificationTestCallbacks) Progress(title string, percent int, resumable bool) {}
func (c *notificationTestCallbacks) WarningSet(warning Warning, message string)         {}
func (c *notificationTestCallbacks) WarningUnset(warning Warning)                       {}
func (c *notificationTestCallbacks) FlushError(message string)                          {}
func (c *notificationTestCallbacks) FatalError(message string)                          {}

func TestNotificationCallbacks(t *testing.T) {
	cb := &notificationTestCallbacks{}
	suite := ChainstateManagerTestSuite{
		MaxBlockHeightToImport: 5,
		NotificationCallbacks:  cb,
	}
	suite.Setup(t)

	if !cb.blockTipCalled {
		t.Error("BlockTip callback was not called")
	}
	if cb.lastBlockHeight != 5 {
		t.Errorf("Expected last block height 5, got %d", cb.lastBlockHeight)
	}
	if !cb.headerTipCalled {
		t.Error("HeaderTip callback was not called")
	}
	if cb.lastHeaderHeight != 5 {
		t.Errorf("Expected last header height 5, got %d", cb.lastHeaderHeight)
	}
}
