// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-06-02 9:36
// version: 1.0.0
// desc   :

package component

import (
	"time"
)

type TickDoer func()

type TimeTicker struct {
	ticker *time.Ticker
	doer   TickDoer
}

func StartTicker(d time.Duration, doer TickDoer) *TimeTicker {
	tt := &TimeTicker{
		ticker: time.NewTicker(d),
		doer:   doer,
	}
	go tt.run()
	return tt
}

func (tt *TimeTicker) Stop() {
	tt.ticker.Stop()
}

func (tt *TimeTicker) run() {
	for {
		select {
		case <-tt.ticker.C:
			tt.doer()
		}
	}
}
