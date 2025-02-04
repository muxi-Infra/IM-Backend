package configs

import "testing"

type MockNotifyer struct {
	addr string
}

func (m *MockNotifyer) Callback(a AppConf) {
	m.addr = a.Cache.Addr
}

func TestAppConf_load(t *testing.T) {
	mocknf := &MockNotifyer{
		addr: "mock_addr",
	}
	ac := &AppConf{
		notifyList: []Notifyer{mocknf},
	}
	data := `
svc:
  - name: "c"
    secret: "chen"
  - name: "y"
    secret: "ye"
db:
  addr: "123567"
  password: "ccnb"
cache:
  addr: "32131"
  password: "ccnb"
  appkeyexpire: 328
  postexpire: 100
  commentexpire: 149
`
	err := ac.load([]byte(data))
	if err != nil {
		t.Error(err)
	}
	t.Log(ac)
	t.Log(mocknf)
	if mocknf.addr != "123567" {
		t.Error("not match")
	}

}
