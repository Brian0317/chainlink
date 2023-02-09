package fluxmonitorv2

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Format implements fmt.Formatter to always print just the pointer address.
// This is a hack to work around a race in github.com/stretchr/testify which
// prints internal fields, including the state of nested, embeded mutexes.
func (fm *FluxMonitor) Format(f fmt.State, verb rune) {
	fmt.Fprintf(f, "%[1]T<%[1]p>", fm)
}

func (fm *FluxMonitor) ExportedPollIfEligible(threshold, absoluteThreshold float64) {
	fm.pollIfEligible(PollRequestTypePoll, NewDeviationChecker(threshold, absoluteThreshold, fm.logger), nil)
}

func (fm *FluxMonitor) ExportedProcessLogs() {
	fm.processLogs()
}

func (fm *FluxMonitor) ExportedBacklog() *utils.BoundedPriorityQueue {
	return fm.backlog
}

func (fm *FluxMonitor) ExportedRoundState() {
	fm.roundState(0)
}

func (fm *FluxMonitor) ExportedRespondToNewRoundLog(log *flux_aggregator_wrapper.FluxAggregatorNewRound, broadcast log.Broadcast) {
	fm.respondToNewRoundLog(*log, broadcast)
}

func (fm *FluxMonitor) ExportedRespondToFlagsRaisedLog() {
	fm.respondToFlagsRaisedLog()
	fm.rotateSelectLoop()
}

func (fm *FluxMonitor) rotateSelectLoop() {
	// the PollRequest is sent to 'rotate' the main select loop, so that new timers will be evaluated
	fm.pollManager.chPoll <- PollRequest{Type: PollRequestTypeUnknown}
}
