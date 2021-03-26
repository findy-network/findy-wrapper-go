package addons

import (
	"fmt"
	"time"

	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/golang/glog"
	"github.com/lainio/err2"
)

const maxTimeout = 30 * time.Second

var delta = 50 * time.Minute

type myClient struct {
	loginTs    time.Time
	*immu
}

type getData struct {
	data
	reply chan data
}

type data struct {
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
		defer err2.CatchTrace(func(err error) {
			glog.Errorln("fatal error in read", err)
		})
		key, value := err2.StrStr.Try(m.immu.Read(get.key))
		get.reply <- data{key, value}
	}
	write := func(set data) {
		defer err2.CatchTrace(func(err error) {
			glog.Errorln("fatal error in write", err)
		})
		err2.Check(m.immu.Write(set.key, set.value))
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
	defer err2.CatchTrace(func(err error) {
		glog.Errorln("fatal error in refresh token", err)
	})

	compareT := m.loginTs.Add(delta)
	timeDiff := time.Until(compareT)
	if timeDiff >= 0 {
		glog.V(3).Infoln("refresh login")
		err2.Check(m.login())
	}
}

func (m *myClient) login() (err error) {
	defer err2.Return(&err)
	err2.Check(m.immu.login())
	m.loginTs = time.Now()
	return nil
}

func (m *myClient) Read(
	key string,
) (
	_ string,
	_ string,
	err error,
) {
	glog.V(100).Infoln("(((( read")
	reply := make(chan data)
	query := getData{
		data: data{
			key: key,
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

func (m *myClient) Write(
	key string,
	value string,
) (
	err error,
) {
	glog.V(100).Infoln("(((( write")
	setChannel <- data{key: key, value: value}
	glog.V(100).Infoln(")))) write")
	return nil
}

func (m *myClient) Close() {
	m.Stop()
	m.immu.Close()
}

func (m *myClient) Open(name string) bool {
	m.Start()
	return m.immu.Open(name)
}

const immuLedgerName = "FINDY_IMMUDB_LEDGER"

var immuLedger = newClient(immuLedgerImpl)

func init() {
	pool.RegisterPlugin(immuLedgerName, immuLedger)
}
