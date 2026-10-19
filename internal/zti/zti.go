package zti

type Status uint8

const (
	StatusSleep Status = iota
	StatusUpdated
	StatusClean
)

type Station uint8

const (
	StationNone Station = iota
	StationZ
	StationT
	StationI
	StationII
	StationCommit
)

type Trace struct {
	Tick    uint64
	Status  Status
	BitZ    uint8 // 1 pass, 0 fail or not reached
	BitT    uint8
	BitI    uint8
	BitII   uint8
	Station Station
}

type Result struct {
	Trace Trace
	Data  any
}

type Chain struct {
	hasher   *Hasher
	registry *Registry
	maxSize  int
	tick     uint64
}

func NewChain() *Chain {
	return &Chain{
		hasher:   NewHasher(8),
		registry: NewRegistry(3),
		maxSize:  3,
	}
}

func begin(c *Chain, t *Trace) {
	c.tick++
	*t = Trace{Tick: c.tick, Status: StatusSleep}
}

func drop(t *Trace, at Station) Result {
	t.Status = StatusClean
	t.Station = at
	return Result{Trace: *t}
}

func commit(t *Trace, data any) Result {
	t.Status = StatusUpdated
	t.Station = StationCommit
	return Result{Trace: *t, Data: data}
}
