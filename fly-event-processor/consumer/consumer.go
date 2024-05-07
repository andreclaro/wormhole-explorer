package consumer

import (
	"context"

	"github.com/wormhole-foundation/wormhole-explorer/common/pool"
	"github.com/wormhole-foundation/wormhole-explorer/fly-event-processor/internal/metrics"
	"github.com/wormhole-foundation/wormhole-explorer/fly-event-processor/processor"
	"github.com/wormhole-foundation/wormhole-explorer/fly-event-processor/queue"
	sdk "github.com/wormhole-foundation/wormhole/sdk/vaa"
	"go.uber.org/zap"
)

// Consumer consumer struct definition.
type Consumer struct {
	consumeFunc  queue.ConsumeFunc
	processor    processor.ProcessorFunc
	guardianPool *pool.Pool
	logger       *zap.Logger
	metrics      metrics.Metrics
	p2pNetwork   string
	workersSize  int
}

// New creates a new vaa consumer.
func New(
	consumeFunc queue.ConsumeFunc,
	processor processor.ProcessorFunc,
	logger *zap.Logger,
	metrics metrics.Metrics,
	p2pNetwork string,
	workersSize int,
) *Consumer {

	c := Consumer{
		consumeFunc: consumeFunc,
		processor:   processor,
		logger:      logger,
		metrics:     metrics,
		p2pNetwork:  p2pNetwork,
		workersSize: workersSize,
	}

	return &c
}

// Start consumes messages from VAA queue, parse and store those messages in a repository.
func (c *Consumer) Start(ctx context.Context) {
	ch := c.consumeFunc(ctx)
	for i := 0; i < c.workersSize; i++ {
		go c.producerLoop(ctx, ch)
	}
}

func (c *Consumer) producerLoop(ctx context.Context, ch <-chan queue.ConsumerMessage) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			c.processEvent(ctx, msg)
		}
	}
}

func (c *Consumer) processEvent(ctx context.Context, msg queue.ConsumerMessage) {
	event := msg.Data()
	vaaID := event.Data.VaaID
	chainID := sdk.ChainID(event.Data.ChainID)

	logger := c.logger.With(
		zap.String("trackId", event.TrackID),
		zap.String("vaaId", vaaID))

	if msg.IsExpired() {
		msg.Failed()
		logger.Debug("event is expired")
		c.metrics.IncDuplicatedVaaExpired(chainID)
		return
	}

	params := &processor.Params{
		TrackID: event.TrackID,
		VaaID:   vaaID,
		ChainID: chainID,
	}

	err := c.processor(ctx, params)
	if err != nil {
		msg.Failed()
		logger.Error("error processing event", zap.Error(err))
		c.metrics.IncDuplicatedVaaFailed(chainID)
		return
	}

	msg.Done()
	logger.Debug("event processed")
	c.metrics.IncDuplicatedVaaProcessed(chainID)
}
