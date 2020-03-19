package parquetgo

import (
	"fmt"
	"github.com/apache/arrow/go/arrow"
	"github.com/repro/columnify/schema"
	"github.com/xitongsys/parquet-go/layout"
	"reflect"
	"testing"
)

func TestMarshalMap(t *testing.T) {
	cases := []struct {
		input  []interface{}
		bgn    int
		end    int
		schema *schema.IntermediateSchema
		expect *map[string]*layout.Table
		err    error
	}{
		{
			input: []interface{}{
				map[string]interface{}{
					"boolean": false,
					"bytes":   fmt.Sprintf("%v", []byte("foo")),
					"double":  1.1,
					"float":   1.1,
					"int":     1,
					"long":    1,
					"string":  "foo",
				},
				map[string]interface{}{
					"boolean": true,
					"bytes":   fmt.Sprintf("%v", []byte("bar")),
					"double":  2.2,
					"float":   2.2,
					"int":     2,
					"long":    2,
					"string":  "bar",
				},
			},
			bgn: 0,
			end: 2,
			schema: schema.NewIntermediateSchema(
				arrow.NewSchema(
					[]arrow.Field{
						{
							Name:     "boolean",
							Type:     arrow.FixedWidthTypes.Boolean,
							Nullable: false,
						},
						{
							Name:     "int",
							Type:     arrow.PrimitiveTypes.Uint32,
							Nullable: false,
						},
						{
							Name:     "long",
							Type:     arrow.PrimitiveTypes.Uint64,
							Nullable: false,
						},
						{
							Name:     "float",
							Type:     arrow.PrimitiveTypes.Float32,
							Nullable: false,
						},
						{
							Name:     "double",
							Type:     arrow.PrimitiveTypes.Float64,
							Nullable: false,
						},
						{
							Name:     "bytes",
							Type:     arrow.BinaryTypes.Binary,
							Nullable: false,
						},
						{
							Name:     "string",
							Type:     arrow.BinaryTypes.String,
							Nullable: false,
						},
					}, nil),
				"primitives"),
			expect: &map[string]*layout.Table{
				"Primitives.Boolean": {
					Values:           []interface{}{false, true},
					DefinitionLevels: []int32{0, 0},
					RepetitionLevels: []int32{0, 0},
				},
				"Primitives.Int": {
					Values:           []interface{}{int32(1), int32(2)},
					DefinitionLevels: []int32{0, 0},
					RepetitionLevels: []int32{0, 0},
				},
				"Primitives.Long": {
					Values:           []interface{}{int64(1), int64(2)},
					DefinitionLevels: []int32{0, 0},
					RepetitionLevels: []int32{0, 0},
				},
				"Primitives.Float": {
					Values:           []interface{}{float32(1.1), float32(2.2)},
					DefinitionLevels: []int32{0, 0},
					RepetitionLevels: []int32{0, 0},
				},
				"Primitives.Double": {
					Values:           []interface{}{float64(1.1), float64(2.2)},
					DefinitionLevels: []int32{0, 0},
					RepetitionLevels: []int32{0, 0},
				},
				"Primitives.Bytes": {
					Values:           []interface{}{fmt.Sprintf("%v", []byte("foo")), fmt.Sprintf("%v", []byte("bar"))},
					DefinitionLevels: []int32{0, 0},
					RepetitionLevels: []int32{0, 0},
				},
				"Primitives.String": {
					Values:           []interface{}{"foo", "bar"},
					DefinitionLevels: []int32{0, 0},
					RepetitionLevels: []int32{0, 0},
				},
			},
			err: nil,
		},
	}

	for _, c := range cases {
		sh, err := schema.NewSchemaHandlerFromArrow(*c.schema)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		tables, err := MarshalMap(c.input, c.bgn, c.end, sh)
		if err != c.err {
			t.Errorf("expected: %v, but actual: %v\n", c.err, err)
		}

		for k, v := range *c.expect {
			actual := (*tables)[k]

			if !reflect.DeepEqual(actual.Values, v.Values) {
				t.Errorf("expected: %v, but actual: %v\n", v.Values, actual.Values)
			}

			if !reflect.DeepEqual(actual.DefinitionLevels, v.DefinitionLevels) {
				t.Errorf("expected: %v, but actual: %v\n", v.DefinitionLevels, actual.DefinitionLevels)
			}

			if !reflect.DeepEqual(actual.RepetitionLevels, v.RepetitionLevels) {
				t.Errorf("expected: %v, but actual: %v\n", v.RepetitionLevels, actual.RepetitionLevels)
			}
		}
	}
}
