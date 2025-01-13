package configs

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewNacosConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want NacosConfig
	}{
		{
			"test1",
			args{path: "config-example.yaml"},
			NacosConfig{
				IpAddr:      "127.0.0.1",
				Port:        8848,
				Namespaceid: "im",
				Group:       "dev",
				DataId:      "app",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNacosConfig(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNacosConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestNacosClient_GetConfig(t *testing.T) {
	nc := NacosConfig{
		IpAddr:      "127.0.0.1",
		Port:        8848,
		Namespaceid: "d26ed8ca-106d-4b49-b2a2-25ac10aa3e07",
		Group:       "dev",
		DataId:      "app",
	}
	client := NewNacosClient(nc)
	data, err := client.GetConfig()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}

func TestNacosClient_ListenConfig(t *testing.T) {
	nc := NacosConfig{
		IpAddr:      "127.0.0.1",
		Port:        8848,
		Namespaceid: "d26ed8ca-106d-4b49-b2a2-25ac10aa3e07",
		Group:       "dev",
		DataId:      "app",
	}
	client := NewNacosClient(nc)
	data, err := client.GetConfig()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
	var wg sync.WaitGroup
	wg.Add(1)
	time.AfterFunc(1*time.Minute, func() {
		wg.Done()
	})
	f := func(data string) {
		defer wg.Done()
		fmt.Println(data)
	}
	if err := client.ListenConfig(f); err != nil {
		t.Error(err)
	}
	wg.Wait()

}
