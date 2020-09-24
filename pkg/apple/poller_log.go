package apple

import (
	"github.com/FTChinese/go-rest/chrono"
	"sync"
)

const StmtSavePollerLog = `
INSERT INTO premium.apple_poller_log
SET total_counter = :total_counter,
	success_counter = :success_counter,
	failure_counter = :failure_counter,
	start_utc = :start_utc,
	end_utc = :end_utc`

type PollerLog struct {
	Total     int64       `db:"total_counter"`
	Succeeded int64       `db:"success_counter"`
	Failed    int64       `db:"failure_counter"`
	StartUTC  chrono.Time `db:"start_utc"`
	EndUTC    chrono.Time `db:"end_utc"`
	mux       sync.Mutex
}

func NewPollerLog() *PollerLog {
	return &PollerLog{
		StartUTC: chrono.TimeNow(),
	}
}

func (p *PollerLog) IncTotal() {
	p.mux.Lock()
	p.Total++
	p.mux.Unlock()
}

func (p *PollerLog) IncSuccess() {
	p.mux.Lock()
	p.Succeeded++
	p.mux.Unlock()
}

func (p *PollerLog) IncFailure() {
	p.mux.Lock()
	p.Failed++
	p.mux.Unlock()
}
