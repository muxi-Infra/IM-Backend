package pkg

import (
	"fmt"
	"reflect"
	"testing"
	"time"
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

func TestAES(t *testing.T) {
	secret := "W7K8pJ3aQv2LcXgH"
	appKey := fmt.Sprintf("%d", time.Now().Unix())
	encryptedKey, err := EncryptAES(appKey, []byte(secret))
	if err != nil {
		t.Error(err)
	}
	t.Log(encryptedKey)
	decryptedKey, err := DecryptAES(encryptedKey, []byte(secret))
	if err != nil {
		t.Error(err)
	}
	if string(decryptedKey) != appKey {
		t.Error("Decrypted key does not match the original key")
	}
}

func TestFormatTimeInShanghai(t *testing.T) {
	// 预加载所需时区
	shanghaiLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("加载上海时区失败: %v", err)
	}
	nyLoc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("加载纽约时区失败: %v", err)
	}
	utcLoc := time.UTC
	type args struct {
		t time.Time
	}
	// 定义测试用例
	tests := []struct {
		name string
		args
		want string
	}{
		{
			name: "上海本地时间",
			args: args{
				t: time.Date(2023, 10, 1, 12, 34, 0, 0, shanghaiLoc),
			},
			want: "2023-10-01T12:34:00",
		},
		{
			name: "UTC时间转换",
			args: args{
				t: time.Date(2023, 10, 1, 4, 30, 0, 0, utcLoc),
			},
			want: "2023-10-01T12:30:00",
		},
		{
			name: "纽约夏令时时间转换",
			args: args{
				t: time.Date(2023, 7, 1, 3, 30, 0, 0, nyLoc), // EDT (UTC-4)
			},
			want: "2023-07-01T15:30:00", // 3:30 EDT → UTC 7:30 → +8小时 = 15:30
		},
		{
			name: "跨天转换",
			args: args{
				t: time.Date(2023, 1, 1, 16, 0, 0, 0, utcLoc), // UTC 16:00 → 上海 00:00 (+1天)
			},
			want: "2023-01-02T00:00:00",
		},
		{
			name: "月末边界时间",
			args: args{
				t: time.Date(2023, 2, 28, 23, 59, 0, 0, utcLoc), // UTC 23:59 → 上海 07:59 (+1天)
			},
			want: "2023-03-01T07:59:00",
		},
		{
			name: "整点零分",
			args: args{
				t: time.Date(2024, 5, 5, 8, 0, 0, 0, utcLoc), // UTC 8:00 → 上海 16:00
			},
			want: "2024-05-05T16:00:00",
		},
		{
			name: "纽约标准时间转换", // 非夏令时（EST UTC-5）
			args: args{
				t: time.Date(2023, 12, 1, 3, 30, 0, 0, nyLoc), // EST 3:30 → UTC 8:30 → 上海 16:30
			},
			want: "2023-12-01T16:30:00",
		},
	}

	// 运行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTimeInShanghai(tt.args.t); got != tt.want {
				t.Errorf("FormatTimeInShanghai() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecryptAES(t *testing.T) {
	appKey := "45483374579c0a4bcafd3b3c68a0336b10357d21013e1420c469d8a7a08e75c2"
	secret := "W7K8pJ3aQv2LcXgH"
	res, err := DecryptAES(appKey, []byte(secret))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Decrypted appKey: %s\n", string(res))
}
