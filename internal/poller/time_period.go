package poller

import "time"

type timePeriod struct {
	start time.Duration
	end   time.Duration
}

func (p *timePeriod) getTimes() (time.Time, time.Time) {
	n := time.Now().UTC()
	return n.Add(p.start), n.Add(p.end)
}
