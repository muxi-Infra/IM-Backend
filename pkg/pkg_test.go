package pkg

import (
	"reflect"
	"testing"
)

func TestReadYamlContent(t *testing.T) {
	type args struct {
		content []byte
		aim     any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid YAML",
			args: args{
				content: []byte("key: value"),
				aim:     &map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "Invalid YAML",
			args: args{
				content: []byte("key: value:"),
				aim:     &map[string]string{},
			},
			wantErr: true,
		},
		{
			name: "Empty content",
			args: args{
				content: []byte(""),
				aim:     &map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "Nested YAML",
			args: args{
				content: []byte("parent:\n  child: value"),
				aim:     &map[string]map[string]string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadYamlContent(tt.args.content, tt.args.aim); (err != nil) != tt.wantErr {
				t.Errorf("ReadYamlContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMergeMaps(t *testing.T) {
	type args struct {
		map1 map[string]int
		map2 map[string]int
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{
			name: "Both maps empty",
			args: args{
				map1: map[string]int{},
				map2: map[string]int{},
			},
			want: map[string]int{},
		},
		{
			name: "First map empty",
			args: args{
				map1: map[string]int{},
				map2: map[string]int{"a": 1, "b": 2},
			},
			want: map[string]int{"a": 1, "b": 2},
		},
		{
			name: "Second map empty",
			args: args{
				map1: map[string]int{"a": 1, "b": 2},
				map2: map[string]int{},
			},
			want: map[string]int{"a": 1, "b": 2},
		},
		{
			name: "No overlapping keys",
			args: args{
				map1: map[string]int{"a": 1},
				map2: map[string]int{"b": 2},
			},
			want: map[string]int{"a": 1, "b": 2},
		},
		{
			name: "Overlapping keys",
			args: args{
				map1: map[string]int{"a": 1, "b": 2},
				map2: map[string]int{"b": 3, "c": 4},
			},
			want: map[string]int{"a": 1, "b": 3, "c": 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeMaps(tt.args.map1, tt.args.map2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeMaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	type args struct {
		slice []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Empty slice",
			args: args{
				slice: []int{},
			},
			want: []int{},
		},
		{
			name: "No duplicates",
			args: args{
				slice: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "With duplicates",
			args: args{
				slice: []int{1, 2, 2, 3, 3, 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "All duplicates",
			args: args{
				slice: []int{1, 1, 1, 1},
			},
			want: []int{1},
		},
		{
			name: "Mixed order",
			args: args{
				slice: []int{3, 1, 2, 3, 2, 1},
			},
			want: []int{3, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}
