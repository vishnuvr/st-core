package core

import "sync"

func ValueStore() SourceSpec {
	return SourceSpec{
		Name: "value",
		Type: VALUE_PRIMITIVE,
		New:  NewValue,
	}
}

type Value struct {
	value interface{}
	quit  chan bool
	sync.Mutex
}

func NewValue() Source {
	return &Value{
		value: nil,
		quit:  make(chan bool),
	}
}

func (v Value) GetType() SourceType {
	return VALUE_PRIMITIVE
}

// ValueGet emits the value stored
func ValueGet() Spec {
	return Spec{
		Name: "valueGet",
		Inputs: []Pin{
			Pin{"trigger"},
		},
		Outputs: []Pin{
			Pin{"value"},
		},
		Source: VALUE_PRIMITIVE,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			v := s.(*Value)
			out[0] = v.value
			return nil
		},
	}
}

// ValueSet sets the value stored
func ValueSet() Spec {
	return Spec{
		Name: "valueSet",
		Inputs: []Pin{
			Pin{"value"},
		},
		Outputs: []Pin{
			Pin{"out"},
		},
		Source: VALUE_PRIMITIVE,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			v := s.(*Value)
			v.value = in[0]
			out[0] = true
			return nil
		},
	}
}
