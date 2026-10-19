package zti

import (
	"reflect"
	"testing"
)

type ipcCase struct {
	name    string
	seed    []IPCEntry
	input   IPCInput
	bitZ    uint8
	bitT    uint8
	bitI    uint8
	bitII   uint8
	station Station
	status  Status
	obj     string
}

func TestStepIPC(t *testing.T) {
	tests := []ipcCase{
		{
			name:  "happy_path",
			seed:  []IPCEntry{{id: 1, obj: "PROC-1", valid: true}},
			input: IPCInput{TargetID: 1, HasID: true},
			bitZ:  1, bitT: 1, bitI: 1, bitII: 1,
			station: StationCommit,
			status:  StatusUpdated,
			obj:     "PROC-1",
		},
		{
			name:    "drop_z_no_int",
			input:   IPCInput{HasID: false},
			bitZ:    0,
			station: StationZ,
			status:  StatusClean,
		},
		{
			name:    "drop_t_empty_bucket",
			input:   IPCInput{TargetID: 1, HasID: true},
			bitZ:    1,
			station: StationT,
			status:  StatusClean,
		},
		{
			name:  "drop_i_id_absent",
			seed:  []IPCEntry{{id: 1, obj: "PROC-1", valid: true}},
			input: IPCInput{TargetID: 9, HasID: true},
			bitZ:  1, bitT: 1,
			station: StationI,
			status:  StatusClean,
		},
		{
			name:  "drop_ii_invalid",
			seed:  []IPCEntry{{id: 1, obj: "PROC-1", valid: false}},
			input: IPCInput{TargetID: 1, HasID: true},
			bitZ:  1, bitT: 1, bitI: 1,
			station: StationII,
			status:  StatusClean,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := NewChain()
			for _, e := range tc.seed {
				c.hasher.Register(e.id, e.obj, e.valid)
			}

			r := c.StepIPC(tc.input)
			tr := r.Trace

			if tr.BitZ != tc.bitZ || tr.BitT != tc.bitT || tr.BitI != tc.bitI || tr.BitII != tc.bitII {
				t.Fatalf("bits = %d%d%d%d, want %d%d%d%d",
					tr.BitZ, tr.BitT, tr.BitI, tr.BitII,
					tc.bitZ, tc.bitT, tc.bitI, tc.bitII)
			}
			if tr.Station != tc.station {
				t.Fatalf("station = %d, want %d", tr.Station, tc.station)
			}
			if tr.Status != tc.status {
				t.Fatalf("status = %d, want %d", tr.Status, tc.status)
			}
			if tc.obj != "" {
				d, ok := r.Data.(IPCData)
				if !ok || d.Obj != tc.obj || d.TargetID != tc.input.TargetID {
					t.Fatalf("data = %#v, want obj=%q id=%d", r.Data, tc.obj, tc.input.TargetID)
				}
			} else if r.Data != nil {
				t.Fatalf("data = %#v, want nil on drop", r.Data)
			}
		})
	}
}

func TestStepIPCDeterminism(t *testing.T) {
	run := func() []Trace {
		c := NewChain()
		c.hasher.Register(1, "PROC-1", true)
		c.hasher.Register(3, "PROC-3", false)
		inputs := []IPCInput{
			{TargetID: 1, HasID: true},
			{TargetID: 2, HasID: true},
			{TargetID: 3, HasID: true},
			{HasID: false},
			{TargetID: 1, HasID: true},
		}
		var out []Trace
		for _, in := range inputs {
			out = append(out, c.StepIPC(in).Trace)
		}
		return out
	}

	a, b := run(), run()
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("nondeterministic traces:\nA=%+v\nB=%+v", a, b)
	}
}

type objectCase struct {
	name    string
	seed    []string
	input   ObjectInput
	bitZ    uint8
	bitT    uint8
	bitI    uint8
	bitII   uint8
	station Station
	status  Status
	handle  string
}

func TestStepObject(t *testing.T) {
	tests := []objectCase{
		{
			name:  "happy_path",
			input: ObjectInput{Command: "REGISTER", Arg: "alpha", HasArg: true},
			bitZ:  1, bitT: 1, bitI: 1, bitII: 1,
			station: StationCommit,
			status:  StatusUpdated,
			handle:  "OBJ-1-1",
		},
		{
			name:    "drop_z_bad_command",
			input:   ObjectInput{Command: "PING", Arg: "alpha", HasArg: true},
			bitZ:    0,
			station: StationZ,
			status:  StatusClean,
		},
		{
			name:    "drop_t_registry_full",
			seed:    []string{"alpha", "beta", "delta"},
			input:   ObjectInput{Command: "REGISTER", Arg: "beta", HasArg: true},
			bitZ:    1,
			station: StationT,
			status:  StatusClean,
		},
		{
			name:  "drop_i_empty_arg",
			input: ObjectInput{Command: "REGISTER", HasArg: true},
			bitZ:  1, bitT: 1,
			station: StationI,
			status:  StatusClean,
		},
		{
			name:  "drop_ii_duplicate_arg",
			seed:  []string{"alpha"},
			input: ObjectInput{Command: "REGISTER", Arg: "alpha", HasArg: true},
			bitZ:  1, bitT: 1, bitI: 1,
			station: StationII,
			status:  StatusClean,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := NewChain()
			for _, arg := range tc.seed {
				c.registry.Register(arg)
			}

			r := c.StepObject(tc.input)
			tr := r.Trace

			if tr.BitZ != tc.bitZ || tr.BitT != tc.bitT || tr.BitI != tc.bitI || tr.BitII != tc.bitII {
				t.Fatalf("bits = %d%d%d%d, want %d%d%d%d",
					tr.BitZ, tr.BitT, tr.BitI, tr.BitII,
					tc.bitZ, tc.bitT, tc.bitI, tc.bitII)
			}
			if tr.Station != tc.station {
				t.Fatalf("station = %d, want %d", tr.Station, tc.station)
			}
			if tr.Status != tc.status {
				t.Fatalf("status = %d, want %d", tr.Status, tc.status)
			}
			if tc.handle != "" {
				d, ok := r.Data.(ObjectData)
				if !ok || d.Handle != tc.handle || d.Arg != tc.input.Arg {
					t.Fatalf("data = %#v, want handle=%q arg=%q", r.Data, tc.handle, tc.input.Arg)
				}
			} else if r.Data != nil {
				t.Fatalf("data = %#v, want nil on drop", r.Data)
			}
		})
	}
}

func TestStepObjectDeterminism(t *testing.T) {
	run := func() []Trace {
		c := NewChain()
		inputs := []ObjectInput{
			{Command: "REGISTER", Arg: "alpha", HasArg: true},
			{Command: "REGISTER", Arg: "beta", HasArg: true},
			{Command: "REGISTER", Arg: "alpha", HasArg: true},
			{Command: "REGISTER", HasArg: true},
			{Command: "PING", Arg: "gamma", HasArg: true},
		}
		var out []Trace
		for _, in := range inputs {
			out = append(out, c.StepObject(in).Trace)
		}
		return out
	}

	a, b := run(), run()
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("nondeterministic traces:\nA=%+v\nB=%+v", a, b)
	}
}
