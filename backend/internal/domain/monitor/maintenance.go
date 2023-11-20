package monitor

import "time"

type Maintenance struct {
	from     time.Time
	to       time.Time
	timezone time.Location
}

func (m Maintenance) From() time.Time {
	return m.from
}

func (m Maintenance) To() time.Time {
	return m.to
}

func (m Maintenance) Timezone() time.Location {
	return m.timezone
}
