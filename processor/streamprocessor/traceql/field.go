package traceql

import (
	"github.com/joe-elliott/trace-streaming/processor/streamprocessor/streampb"
)

const (
	fieldTypeUnknown = -1
	fieldTypeInt     = 1
	fieldTypeFloat   = 2
	fieldTypeString  = 3
)

type fieldID []int

type field interface {
	getIntValue(*streampb.Span) int
	getFloatValue(*streampb.Span) float64
	getStringValue(*streampb.Span) string

	getNativeType(*streampb.Span) int
	getRelationshipID() fieldID
}

type intField int

func newIntField(val int) field {
	return intField(val)
}

func (n intField) getIntValue(span *streampb.Span) int {
	return int(n)
}

func (n intField) getFloatValue(span *streampb.Span) float64 {
	return float64(n)
}

func (n intField) getStringValue(span *streampb.Span) string {
	return ""
}

func (n intField) getNativeType(span *streampb.Span) int {
	return fieldTypeInt
}

func (n intField) getRelationshipID() fieldID {
	return fieldID(nil)
}

type floatField float64

func newFloatField(val float64) field {
	return floatField(val)
}

func (f floatField) getIntValue(span *streampb.Span) int {
	return 0
}

func (f floatField) getFloatValue(span *streampb.Span) float64 {
	return float64(f)
}

func (f floatField) getStringValue(span *streampb.Span) string {
	return ""
}

func (f floatField) getNativeType(span *streampb.Span) int {
	return fieldTypeFloat
}

func (f floatField) getRelationshipID() fieldID {
	return fieldID(nil)
}

type stringField string

func newStringField(val string) field {
	return stringField(val)
}

func (s stringField) getIntValue(span *streampb.Span) int {
	return 0
}

func (s stringField) getFloatValue(span *streampb.Span) float64 {
	return 0.0
}

func (s stringField) getStringValue(span *streampb.Span) string {
	return string(s)
}

func (s stringField) getNativeType(span *streampb.Span) int {
	return fieldTypeString
}

func (s stringField) getRelationshipID() fieldID {
	return fieldID(nil)
}

// dynamicField
type dynamicField struct {
	relID fieldID
	id    fieldID
	name  string // only valid if id = FIELD_ATTS or FIELD_EVENTS
}

func newDynamicField(id int, name string) field {
	return dynamicField{
		id:   []int{id},
		name: name,
	}
}

func wrapDynamicField(id int, f field) field {
	d := f.(dynamicField)

	return dynamicField{
		relID: d.relID,
		id:    append([]int{id}, d.id...),
		name:  d.name,
	}
}

func wrapRelationshipField(id int, f field) field {
	d := f.(dynamicField)

	return dynamicField{
		relID: append([]int{id}, d.relID...),
		id:    d.id,
		name:  d.name,
	}
}

func (f dynamicField) getIntValue(s *streampb.Span) int {
	if len(f.id) == 0 {
		return 0
	}

	rootID := f.id[0]

	switch rootID {
	case FIELD_DURATION:
		return int(s.Duration)
	case FIELD_ATTS:
		if a, ok := s.Attributes[f.name]; ok {
			switch a.Type {
			case streampb.KeyValuePair_INT:
				return int(a.IntValue)
			case streampb.KeyValuePair_BOOL:
				if a.BoolValue {
					return 1
				}

				return 0
			}
		}
	case FIELD_EVENTS:
		if e, ok := s.Events[f.name]; ok {
			switch e.Type {
			case streampb.KeyValuePair_INT:
				return int(e.IntValue)
			case streampb.KeyValuePair_BOOL:
				if e.BoolValue {
					return 1
				}

				return 0
			}
		}
	case FIELD_STATUS:
		// unsafe check for code/msg
		subfield := f.id[1]
		if subfield == FIELD_CODE {
			return int(s.Status.Code)
		}
	case FIELD_IS_ROOT:
		isRoot := 0
		if len(s.ParentSpanID) == 0 {
			isRoot = 1
		}
		return isRoot
	}

	return 0
}

func (f dynamicField) getFloatValue(s *streampb.Span) float64 {
	if len(f.id) == 0 {
		return 0.0
	}

	rootID := f.id[0]

	switch rootID {
	case FIELD_DURATION:
		return float64(s.Duration)

	case FIELD_ATTS:
		if a, ok := s.Attributes[f.name]; ok {
			switch a.Type {
			case streampb.KeyValuePair_INT:
				return float64(a.IntValue)
			case streampb.KeyValuePair_DOUBLE:
				return a.DoubleValue
			}
		}
	case FIELD_EVENTS:
		if e, ok := s.Events[f.name]; ok {
			switch e.Type {
			case streampb.KeyValuePair_INT:
				return float64(e.IntValue)
			case streampb.KeyValuePair_DOUBLE:
				return e.DoubleValue
			}
		}
	}

	return 0.0
}

func (f dynamicField) getStringValue(s *streampb.Span) string {
	if len(f.id) == 0 {
		return ""
	}

	rootID := f.id[0]

	switch rootID {
	case FIELD_NAME:
		return s.Name
	case FIELD_ATTS:
		if a, ok := s.Attributes[f.name]; ok {
			if a.Type == streampb.KeyValuePair_STRING {
				return a.StringValue
			}
		}
	case FIELD_EVENTS:
		if e, ok := s.Events[f.name]; ok {
			if e.Type == streampb.KeyValuePair_STRING {
				return e.StringValue
			}
		}
	case FIELD_STATUS:
		// unsafe check for code/msg
		subfield := f.id[1]
		if subfield == FIELD_MSG {
			return s.Status.Message
		}
	case FIELD_PROCESS:
		// unsafe check
		subfield := f.id[1]
		if subfield == FIELD_NAME {
			return s.Process.Name
		}
	}

	return ""
}

func (f dynamicField) getNativeType(s *streampb.Span) int {
	if len(f.id) == 0 {
		return fieldTypeUnknown
	}

	rootID := f.id[0]

	switch rootID {
	case FIELD_DURATION:
		return fieldTypeInt
	case FIELD_NAME:
		return fieldTypeString
	case FIELD_ATTS:
		if a, ok := s.Attributes[f.name]; ok {
			switch a.Type {
			case streampb.KeyValuePair_DOUBLE:
				return fieldTypeFloat
			case streampb.KeyValuePair_INT:
				return fieldTypeInt
			case streampb.KeyValuePair_STRING:
				return fieldTypeString
			case streampb.KeyValuePair_BOOL:
				return fieldTypeInt
			}
		}
		return fieldTypeUnknown
	case FIELD_EVENTS:
		if e, ok := s.Events[f.name]; ok {
			switch e.Type {
			case streampb.KeyValuePair_DOUBLE:
				return fieldTypeFloat
			case streampb.KeyValuePair_INT:
				return fieldTypeInt
			case streampb.KeyValuePair_STRING:
				return fieldTypeString
			case streampb.KeyValuePair_BOOL:
				return fieldTypeInt
			}
		}
		return fieldTypeUnknown
	case FIELD_STATUS:
		// unsafe check for code/msg
		subfield := f.id[1]
		if subfield == FIELD_CODE {
			return fieldTypeInt
		}
		if subfield == FIELD_MSG {
			return fieldTypeString
		}
		return fieldTypeUnknown
	case FIELD_PROCESS:
		// unsafe check
		subfield := f.id[1]
		if subfield == FIELD_NAME {
			return fieldTypeString
		}
		return fieldTypeUnknown
	case FIELD_IS_ROOT:
		return fieldTypeInt
	}

	return fieldTypeUnknown
}

func (f dynamicField) getRelationshipID() fieldID {
	return f.relID
}
