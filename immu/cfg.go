package immu

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"time"

	im "github.com/codenotary/immudb/pkg/client"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
)

const (
	// immuMockLedgerName = "FINDY_MOCK_IMMUDB_LEDGER"

	mockURL = "mock"

	envImmuURL  = "ImmuUrl"
	envImmuPort = "ImmuPort"
	envImmuUser = "ImmuUsrName"
	envImmuPwd  = "ImmuPasswd"
)

type Cfg struct {
	URL      string `json:"url"`
	Port     int    `json:"port"`
	UserName string `json:"user_name"`
	Password string `json:"password"`

	*im.Options
}

var MockCfg = &Cfg{
	URL:      "mock",
	Port:     3322,
	UserName: "immudb",
	Password: "immudb",
}

func NewImmuCfg(_ string) (cfg *Cfg) {
	if envExists(envImmuURL) {
		glog.V(2).Infoln("+++ using env Cfg")
		cfg = cfgFromEnv()
	} else {
		glog.V(2).Infoln("+++ using MockCfg")
		cfg = MockCfg
	}
	cfg.Options = im.DefaultOptions().
		WithAddress(cfg.URL).
		WithPort(cfg.Port).
		WithAuth(true)
	return cfg
}

func cfgFromEnv() (cfg *Cfg) {
	assert.D.True(envExists(envImmuURL), "immu URL must exist")
	assert.D.True(envExists(envImmuPort), "immu port must exist")
	assert.D.True(envExists(envImmuUser), "immu user name must exists")
	assert.D.True(envExists(envImmuPwd), "immu password must exist")

	tmpCfg := &Cfg{
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

func (cfg *Cfg) Connect() (c im.ImmuClient, token string, err error) {
	defer err2.Return(&err)

	client, err := cfg.newImmuClient()
	try.To(err)

	token, err = cfg.login(client)
	try.To(err)

	createTokenDir() // for immuDB bug, to allow Logout()
	return client, token, nil
}

func (cfg *Cfg) login(client im.ImmuClient) (token string, err error) {
	defer err2.Return(&err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	lr, err := client.Login(ctx, []byte(cfg.UserName), []byte(cfg.Password))
	try.To(err)

	return lr.Token, nil
}

// createTokenDir because of the bug in immuDB Logout() which cannot be called
// if this empty directory doesn't exist! The immuDB code first tries to delete
// the directory and if it doesn't success it doesn't perform the actual logout.
func createTokenDir() {
	hd, _ := os.UserHomeDir()
	fp := filepath.Join(hd, "token")
	err := os.Mkdir(fp, 0775)
	if err != nil {
		glog.Errorln("mkdir error:", err)
	}
	glog.V(12).Infoln("token path created:", fp)
}

func (cfg *Cfg) newImmuClient() (c im.ImmuClient, err error) {
	if cfg.URL == mockURL {
		return &mockImmuClient{}, nil
	}
	return im.NewImmuClient(cfg.Options)
}
