package helpers

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	_ "github.com/findy-network/findy-wrapper-go/addons/mem" // we need this here

	"github.com/findy-network/findy-wrapper-go/pool"
	"github.com/findy-network/findy-wrapper-go/wallet"
)

func createPool(t *testing.T, maxTimeout time.Duration, poolName string) {
	err := os.RemoveAll(os.Getenv("HOME") + "/.indy_client")
	if err != nil {
		t.Log(err.Error())
	}
	currentPath, _ := os.Getwd()
	genesisPath := currentPath + "/../.circleci/genesis_transactions"
	config := pool.Config{
		GenesisTxn: genesisPath,
	}
	t.Log("Genesis path " + genesisPath)
	select {
	case r := <-pool.CreateConfig(poolName, config):
		if r.Err() != nil {
			t.Log("Pool create error, already exists?")
		}
	case <-time.After(maxTimeout):
		t.Fatal("Timeout exceeded")
	}
}

// OpenTestPool is a helper function for tests to open a ledger pool.
func OpenTestPool(t *testing.T) int {
	r := <-pool.SetProtocolVersion(2)
	if r.Err() != nil {
		t.Fatal("Cannot set pool protocol version")
	}

	const maxTimeout = 5 * time.Second
	poolName := os.Getenv("FINDY_POOL")
	if poolName == "" {
		poolName = "FINDY_MEM_LEDGER"
	}

	if os.Getenv("CI") == "true" {
		createPool(t, maxTimeout, "myNewPool")
	}

	select {
	case r = <-pool.OpenLedger(poolName):
		if r.Err() != nil {
			t.Fatal("Cannot open pool")
		}
	case <-time.After(maxTimeout):
		t.Fatal("Timeout exceeded")
	}

	return r.Handle()
}

//func CreateAndOpenTestWallet2(t *testing.T) int {
//	ut := time.Now().Unix() - 1558885840
//	walletName := fmt.Sprintf("test2_wallet_%v", ut)
//
//	r := <-wallet.Create(wallet.Config{ID: walletName}, wallet.Credentials{Key: "C7mR5TZVB7WRCYsTMQGXuLHcXisFYZL1GoXARyiVyEER", KeyDerivationMethod: "RAW"})
//	if r.Err() != nil {
//		t.Fatal("Cannot create test wallet")
//	}
//	r = <-wallet.Open(wallet.Config{ID: walletName}, wallet.Credentials{Key: "C7mR5TZVB7WRCYsTMQGXuLHcXisFYZL1GoXARyiVyEER", KeyDerivationMethod: "RAW"})
//	if r.Err() != nil {
//		t.Fatal("Cannot open test wallet")
//	}
//	return r.Handle()
//}

var nameCounter uint32
var lock sync.Mutex

func walletName() string {
	lock.Lock()
	defer lock.Unlock()
	nameCounter++
	return fmt.Sprintf("test_wallet_%v_%v", os.Getpid(), nameCounter)
}

// CreateAndOpenTestWallet is a helper function for tests to create and open a
// new wallet. It returns a wallet handle and a wallet name. It generates unique
// name for wallets based on time. Use CloseAndDeleteTestWallet for cleaning up.
func CreateAndOpenTestWallet(t *testing.T) (handle int, name string) {
	walletName := walletName()

	r := <-wallet.Create(wallet.Config{ID: walletName}, wallet.Credentials{Key: "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp", KeyDerivationMethod: "RAW"})
	if r.Err() != nil {
		t.Fatal("Cannot create test wallet")
	}
	r = <-wallet.Open(wallet.Config{ID: walletName}, wallet.Credentials{Key: "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp", KeyDerivationMethod: "RAW"})
	if r.Err() != nil {
		t.Fatal("Cannot open test wallet")
	}
	return r.Handle(), walletName
}

// CloseAndDeleteTestWallet is a helper function for tests to close and delete a
// test wallet.
func CloseAndDeleteTestWallet(w int, name string, t *testing.T) {
	r := <-wallet.Close(w)
	if r.Err() != nil {
		t.Error("Cannot close test wallet")
	}
	r = <-wallet.Delete(wallet.Config{ID: name}, wallet.Credentials{Key: "6cih1cVgRH8yHD54nEYyPKLmdv67o8QbufxaTHot3Qxp", KeyDerivationMethod: "RAW"})
	if r.Err() != nil {
		t.Error("Cannot Delete test wallet")
	}
}

//
//func CloseAndDeleteTestWallet2(w int, t *testing.T) {
//	r := <-wallet.Close(w)
//	if r.Err() != nil {
//		t.Error("Cannot close test wallet2")
//	}
//	r = <-wallet.Delete(wallet.Config{ID: "test2_wallet"}, wallet.Credentials{Key: "C7mR5TZVB7WRCYsTMQGXuLHcXisFYZL1GoXARyiVyEER", KeyDerivationMethod: "RAW"})
//	if r.Err() != nil {
//		t.Error("Cannot Delete test wallet2")
//	}
//}

// CloseTestPool is a helper function for tests. It closes ledger pool used for
// testing. It doesn't do any cleanup.
func CloseTestPool(p int, t *testing.T) {
	r := <-pool.CloseLedger(p)
	if r.Err() != nil {
		t.Error("Cannot close pool")
	}
}
