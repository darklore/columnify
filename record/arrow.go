package record

import (
	"fmt"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
)

type WrappedRecord struct {
	Record array.Record
}

func NewWrappedRecord(b *array.RecordBuilder) *WrappedRecord {
	return &WrappedRecord{
		Record: b.NewRecord(),
	}
}

func formatMapToArrowRecord(s *arrow.Schema, maps []map[string]interface{}) (*WrappedRecord, error) {
	pool := memory.NewGoAllocator()
	b := array.NewRecordBuilder(pool, s)
	defer b.Release()

	for _, m := range maps {
		for i, f := range s.Fields() {
			if v, ok := m[f.Name]; ok {
				if _, err := formatMapToArrowField(b.Field(i), f.Type, v); err != nil {
					return nil, err
				}
			} else {
				b.Field(i).AppendNull()
			}
		}
	}

	return NewWrappedRecord(b), nil
}

func formatMapToArrowStruct(b *array.StructBuilder, s *arrow.StructType, m map[string]interface{}) (*array.StructBuilder, error) {
	for i, f := range s.Fields() {
		if v, ok := m[f.Name]; ok {
			if _, err := formatMapToArrowField(b.FieldBuilder(i), f.Type, v); err != nil {
				return nil, err
			}
		} else {
			b.FieldBuilder(i).AppendNull()
		}

	}

	return b, nil
}

func formatMapToArrowList(b *array.ListBuilder, l *arrow.ListType, list []interface{}) (*array.ListBuilder, error) {
	for _, e := range list {
		if _, err := formatMapToArrowField(b.ValueBuilder(), l.Elem(), e); err != nil {
			return nil, err
		}
	}

	return b, nil
}

func formatMapToArrowField(b array.Builder, t arrow.DataType, v interface{}) (array.Builder, error) {
	switch t.ID() {
	case arrow.BOOL:
		vb, builderOk := b.(*array.BooleanBuilder)
		vv, valueOk := v.(bool)
		if builderOk && valueOk {
			vb.Append(vv)
		} else {
			return nil, fmt.Errorf("unexpected input %v as bool; %w", v, ErrUnconvertibleRecord)
		}

	case arrow.UINT32:
		vb, builderOk := b.(*array.Uint32Builder)
		if !builderOk {
			return nil, fmt.Errorf("builder %v is wrong; %w", v, ErrUnconvertibleRecord)
		}
		switch vv := v.(type) {
		case int32:
			vb.Append(uint32(vv))
		case float64:
			vb.Append(uint32(vv))
		default:
			return nil, fmt.Errorf("unexpected input %v as uint32; %w", v, ErrUnconvertibleRecord)
		}

	case arrow.UINT64:
		vb, builderOk := b.(*array.Uint64Builder)
		if !builderOk {
			return nil, fmt.Errorf("builder %v is wrong; %w", v, ErrUnconvertibleRecord)
		}
		switch vv := v.(type) {
		case int64:
			vb.Append(uint64(vv))
		case float64:
			vb.Append(uint64(vv))
		default:
			return nil, fmt.Errorf("unexpected input %v as uint64; %w", v, ErrUnconvertibleRecord)
		}

	case arrow.FLOAT32:
		vb, builderOk := b.(*array.Float32Builder)
		if !builderOk {
			return nil, fmt.Errorf("builder %v is wrong; %w", v, ErrUnconvertibleRecord)
		}
		switch vv := v.(type) {
		case float32:
			vb.Append(float32(vv))
		case float64:
			vb.Append(float32(vv))
		default:
			return nil, fmt.Errorf("unexpected input %v as float32; %w", v, ErrUnconvertibleRecord)
		}

	case arrow.FLOAT64:
		vb, builderOk := b.(*array.Float64Builder)
		vv, valueOk := v.(float64)
		if builderOk && valueOk {
			vb.Append(vv)
		} else {
			return nil, fmt.Errorf("unexpected input %v as float64; %w", v, ErrUnconvertibleRecord)
		}

	case arrow.STRING:
		vb, builderOk := b.(*array.StringBuilder)
		vv, valueOk := v.(string)
		if builderOk && valueOk {
			vb.Append(vv)
		} else {
			return nil, fmt.Errorf("unexpected input %v as string; %w", v, ErrUnconvertibleRecord)
		}

	case arrow.BINARY:
		vb, builderOk := b.(*array.BinaryBuilder)
		if !builderOk {
			return nil, fmt.Errorf("builder %v is wrong; %w", v, ErrUnconvertibleRecord)
		}
		switch vv := v.(type) {
		case string:
			vb.Append([]byte(vv))
		case []byte:
			vb.Append(vv)
		default:
			return nil, fmt.Errorf("unexpected input %v as binary; %w", v, ErrUnconvertibleRecord)
		}

	case arrow.STRUCT:
		vb, builderOk := b.(*array.StructBuilder)
		st, structOk := t.(*arrow.StructType)
		if builderOk && structOk {
			if v != nil {
				vb.Append(true)
				vv, valueOk := v.(map[string]interface{})
				if !valueOk {
					return nil, fmt.Errorf("unexpected input %v as struct; %w", v, ErrUnconvertibleRecord)
				} else if _, err := formatMapToArrowStruct(vb, st, vv); err != nil {
					return nil, err
				}
			} else {
				vb.Append(false)
			}
		} else {
			return nil, fmt.Errorf("unexpected input %v as struct; %w", v, ErrUnconvertibleRecord)
		}

	case arrow.LIST:
		vb, builderOk := b.(*array.ListBuilder)
		lt, listOk := t.(*arrow.ListType)
		if builderOk && listOk {
			if v != nil {
				vb.Append(true)
				vv, valueOk := v.([]interface{})
				if !valueOk {
					return nil, fmt.Errorf("unexpected input %v as list; %w", v, ErrUnconvertibleRecord)
				}
				if _, err := formatMapToArrowList(vb, lt, vv); err != nil {
					return nil, err
				}
			} else {
				vb.Append(false)
			}
		} else {
			return nil, fmt.Errorf("unexpected input %v as list; %w", v, ErrUnconvertibleRecord)
		}

	default:
		return nil, fmt.Errorf("unconvertable type %v; %w", t.ID(), ErrUnconvertibleRecord)
	}

	return b, nil
}
