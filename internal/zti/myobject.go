package zti

import "strconv"

type Registry struct {
	slots map[int]string
	free  []int
	next  int
	max   int
}

func NewRegistry(max int) *Registry {
	return &Registry{slots: make(map[int]string), next: 1, max: max}
}

func (r *Registry) alloc() (int, bool) {
	if n := len(r.free); n > 0 {
		slot := r.free[n-1]
		r.free = r.free[:n-1]
		return slot, true
	}
	if r.next <= r.max {
		slot := r.next
		r.next++
		return slot, true
	}
	return 0, false
}

func (r *Registry) release(slot int) {
	r.free = append(r.free, slot)
}

func (r *Registry) Register(arg string) string {
	slot, ok := r.alloc()
	if !ok {
		return ""
	}
	r.slots[slot] = arg
	return "OBJ-" + strconv.Itoa(slot) + "-1"
}

type ObjectInput struct {
	Command string
	Arg     string
	HasArg  bool
}

type ObjectData struct {
	Handle string
	Arg    string
}

func (c *Chain) StepObject(raw ObjectInput) Result {
	var t Trace
	begin(c, &t)

	if raw.Command != "REGISTER" || !raw.HasArg {
		return drop(&t, StationZ)
	}
	t.BitZ = 1

	if len(c.registry.slots) >= c.maxSize {
		return drop(&t, StationT)
	}
	t.BitT = 1

	if len(raw.Arg) == 0 {
		return drop(&t, StationI)
	}
	t.BitI = 1

	for _, v := range c.registry.slots {
		if v == raw.Arg {
			return drop(&t, StationII)
		}
	}
	t.BitII = 1

	handle := c.registry.Register(raw.Arg)
	return commit(&t, ObjectData{Handle: handle, Arg: raw.Arg})
}

func RunObject() *Result {
	c := NewChain()
	r := c.StepObject(ObjectInput{Command: "REGISTER", Arg: "alpha", HasArg: true})
	return &r
}
