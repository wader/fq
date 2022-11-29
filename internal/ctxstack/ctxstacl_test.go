package ctxstack_test

import (
	"testing"

	"github.com/wader/fq/internal/ctxstack"
)

func TestCancelBeforePush(t *testing.T) {
	// TODO: nicer way to test trigger before any push
	waitTriggerFn := make(chan struct{})
	triggerCh := make(chan struct{})
	waitCh := make(chan struct{})
	hasTriggeredOnce := false

	ctxstack.New(func(stopCh chan struct{}) {
		if hasTriggeredOnce {
			close(stopCh)
			close(waitCh)
			return
		}

		close(waitTriggerFn)
		<-triggerCh

		hasTriggeredOnce = true
	})

	// wait for trigger func to be called
	<-waitTriggerFn
	// make trigger func return and cancel
	close(triggerCh)

	<-waitCh
}
