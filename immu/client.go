package immu

import (
	"fmt"
	"time"

	"github.com/findy-network/findy-wrapper-go/plugin"
	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/golang/glog"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
)

const maxTimeout = 30 * time.Second

var delta = 50 * time.Minute

type myClient struct {
	loginTS time.Time
	*immu
}

type getData struct {
	data
	reply chan data
}

type data struct {
	plugin.TxInfo
	key   string
	value string
}

var (
	queryChannel = make(chan getData)
	setChannel   = make(chan data)
	stopChannel  = make(chan struct{})
)

func newClient(immuLedger *immu) *myClient {
	return &myClient{immu: immuLedger}
}

func (m *myClient) Start() {
	// code block for version with separated tokens for reading and writing
	//	go func() {
	//		for get := range queryChannel {
	//			read(getData)
	//		}
	//	}()
	//	go func() {
	//		for set := range setChannel {
	//			write(setData)
	//		}
	//	}()

	read := func(get getData) {
		defer err2.Catch(err2.Err(func(err error) {
			glog.Errorln("fatal error in read", err)
		}))
		key, value := try.To2(m.immu.Read(get.TxInfo, get.key))
		get.reply <- data{get.TxInfo, key, value}
	}
	write := func(set data) {
		defer err2.Catch(err2.Err(func(err error) {
			glog.Errorln("fatal error in write", err)
		}))
		try.To(m.immu.Write(set.TxInfo, set.key, set.value))
	}

	go func() {
		glog.V(1).Infoln("starting immu plugin")
	loop:
		for {
			select {
			case readData := <-queryChannel:
				m.refreshToken()
				read(readData)
			case writeData := <-setChannel:
				m.refreshToken()
				write(writeData)
			case <-stopChannel:
				glog.V(1).Infoln("stopping immu plugin")
				break loop
			}
		}
	}()
}

func (m *myClient) Stop() {
	// spool version
	//	close(queryChannel)
	//	close(setChannel)
	stopChannel <- struct{}{}
}

func (m *myClient) refreshToken() {
	defer err2.Catch(err2.Err(func(err error) {
		glog.Errorln("fatal error in refresh token", err)
	}))

	if m.needRefresh() {
		glog.V(3).Infoln("refresh login")
		try.To(m.login())
	}
}

func (m *myClient) needRefresh() bool {
	expTime := m.loginTS.Add(delta)
	timeDiff := time.Until(expTime)
	return timeDiff <= 0 // no time left i.e. over time
}

func (m *myClient) login() (err error) {
	defer err2.Handle(&err)
	glog.V(1).Infoln("++ login")
	try.To(m.immu.login())
	m.loginTS = time.Now()
	return nil
}

func (m *myClient) Read(tx plugin.TxInfo, key string) (_ string, _ string, err error) {
	glog.V(100).Infoln("(((( read")
	reply := make(chan data)
	query := getData{
		data: data{
			TxInfo: tx,
			key:    key,
		},
		reply: reply,
	}
	queryChannel <- query
	select {
	case r := <-reply:
		glog.V(100).Infoln(")))) read")
		return r.key, r.value, nil
	case <-time.After(maxTimeout):
		return "", "", fmt.Errorf("timeout error")
	}
}

func (m *myClient) Write(tx plugin.TxInfo, key string, value string) (err error) {
	glog.V(100).Infoln("(((( write")
	setChannel <- data{TxInfo: tx, key: key, value: value}
	glog.V(100).Infoln(")))) write")
	return nil
}

func (m *myClient) Close() {
	m.Stop()
	m.immu.Close()
}

func (m *myClient) Open(name ...string) bool {
	m.Start()
	m.loginTS = time.Now() // set it here because Open does the 1st login
	return m.immu.Open(name[0])
}

const immuLedgerName = "FINDY_IMMUDB_LEDGER"

var immuLedger = newClient(immuLedgerImpl)

func init() {
	pool.RegisterPlugin(immuLedgerName, immuLedger)
}
