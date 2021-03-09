package addons

import (
	"context"
	"os"
	"strconv"
	"time"

	immuclient "github.com/codenotary/immudb/pkg/client"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
)

const (
	immuMockLedgerName = "FINDY_MOCK_IMMUDB_LEDGER"

	mockURL = "mock"

	envImmuURL  = "ImmuUrl"
	envImmuPort = "ImmuPort"
	envImmuUser = "ImmuUsrName"
	envImmuPwd  = "ImmuPasswd"
)

type ImmuCfg struct {
	URL      string `json:"url"`
	Port     int    `json:"port"`
	UserName string `json:"user_name"`
	Password string `json:"password"`

	immuclient.Options
}

var MockCfg = &ImmuCfg{
	URL:      "mock",
	Port:     3322,
	UserName: "immudb",
	Password: "immudb",
}

func NewImmuCfg(name string) (cfg *ImmuCfg) {
	if name != immuMockLedgerName && envExists(envImmuURL) {
		cfg = cfgFromEnv()
	} else {
		glog.V(2).Infoln("using MockCfg")
		cfg = MockCfg
	}
	cfg.Options = immuclient.Options{
		Address:         cfg.URL,
		Port:            cfg.Port,
		Auth:            true,
		CurrentDatabase: "defaultdb",
	}
	return cfg
}

func cfgFromEnv() (cfg *ImmuCfg) {
	assert.D.True(envExists(envImmuURL), "immu URL must exist")
	assert.D.True(envExists(envImmuPort), "immu port must exist")
	assert.D.True(envExists(envImmuUser), "immu user name must exists")
	assert.D.True(envExists(envImmuPwd), "immu password must exist")

	tmpCfg := &ImmuCfg{
		URL:      os.Getenv(envImmuURL),
		Port:     err2.Int.Try(strconv.Atoi(os.Getenv(envImmuPort))),
		UserName: os.Getenv(envImmuUser),
		Password: os.Getenv(envImmuPwd),
	}

	assert.D.True(tmpCfg.URL != "", "immu URL cannot be empty")
	assert.D.True(tmpCfg.Port != 0, "immu port cannot be 0")
	assert.D.True(tmpCfg.UserName != "", "immu user name cannot be empty")
	assert.D.True(tmpCfg.Password != "", "immu password cannot be empty")

	return tmpCfg
}

func envExists(name string) bool {
	_, exists := os.LookupEnv(name)
	return exists
}

func (cfg *ImmuCfg) Connect() (c immuclient.ImmuClient, token string, err error) {
	defer err2.Return(&err)

	client, err := cfg.newImmuClient()
	err2.Check(err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	lr, err := client.Login(ctx, []byte(cfg.UserName), []byte(cfg.Password))
	err2.Check(err)

	return client, lr.Token, nil
}

func (cfg *ImmuCfg) newImmuClient() (c immuclient.ImmuClient, err error) {
	if cfg.URL == mockURL {
		return &mockImmuClient{}, nil
	}
	return immuclient.NewImmuClient(&cfg.Options)
}
