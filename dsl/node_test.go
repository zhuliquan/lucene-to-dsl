package dsl

import (
	"reflect"
	"testing"
)

func TestOrDSLNode_UnionJoin(t *testing.T) {
	type args struct {
		node DSLNode
	}
	tests := []struct {
		name    string
		n       *OrDSLNode
		args    args
		want    DSLNode
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.UnionJoin(tt.args.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrDSLNode.UnionJoin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrDSLNode.UnionJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}
