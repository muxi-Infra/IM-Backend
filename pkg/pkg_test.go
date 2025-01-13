package pkg

import "testing"

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
