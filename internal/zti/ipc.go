package zti

type IPCEntry struct {
	id    int
	obj   string
	valid bool
}

type Hasher struct {
	buckets [][]IPCEntry
	n       int
}

func NewHasher(n int) *Hasher {
	if n <= 0 {
		n = 8
	}
	return &Hasher{buckets: make([][]IPCEntry, n), n: n}
}

func (h *Hasher) bucket(id int) int {
	return id % h.n
}

func (h *Hasher) Register(id int, obj string, valid bool) {
	b := h.bucket(id)
	h.buckets[b] = append(h.buckets[b], IPCEntry{id: id, obj: obj, valid: valid})
}

type IPCInput struct {
	TargetID int
	HasID    bool
}

type IPCData struct {
	TargetID int
	Obj      string
}

func (c *Chain) StepIPC(raw IPCInput) Result {
	var t Trace
	begin(c, &t)

	if !raw.HasID {
		return drop(&t, StationZ)
	}
	t.BitZ = 1

	bucket := c.hasher.buckets[c.hasher.bucket(raw.TargetID)]
	if len(bucket) == 0 {
		return drop(&t, StationT)
	}
	t.BitT = 1

	idx := -1
	for i := range bucket {
		if bucket[i].id == raw.TargetID {
			idx = i
			break
		}
	}
	if idx < 0 {
		return drop(&t, StationI)
	}
	t.BitI = 1

	if !bucket[idx].valid {
		return drop(&t, StationII)
	}
	t.BitII = 1

	return commit(&t, IPCData{TargetID: raw.TargetID, Obj: bucket[idx].obj})
}

func RunIPC() *Result {
	c := NewChain()
	c.hasher.Register(1, "PROC-1", true)
	r := c.StepIPC(IPCInput{TargetID: 1, HasID: true})
	return &r
}
