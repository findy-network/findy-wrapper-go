package pool_test

import (
	"reflect"
	"testing"

	_ "github.com/optechlab/findy-go/addons/echo"
	_ "github.com/optechlab/findy-go/addons/mem"
	"github.com/optechlab/findy-go/pool"
	"github.com/stretchr/testify/assert"
)

func TestSetProtocolVersion(t *testing.T) {
	type args struct {
		version uint64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Protocol version 1", args{1}, 0},
		{"Protocol version 2", args{2}, 0},
		{"Unsupported Protocol version 0", args{3}, 308},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (<-pool.SetProtocolVersion(tt.args.version)).ErrCode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetProtocolVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenLedger(t *testing.T) {
	r := <-pool.OpenLedger("FINDY_MEM_LEDGER", "FINDY_ECHO_LEDGER")
	assert.NoError(t, r.Err())
	h1 := r.Handle()
	assert.Equal(t, h1, -2)

	r = <-pool.CloseLedger(h1)
	assert.NoError(t, r.Err())
}

func TestCloseLedger(t *testing.T) {
	r := <-pool.OpenLedger("FINDY_MEM_LEDGER", "FINDY_ECHO_LEDGER")
	assert.NoError(t, r.Err())
	h1 := r.Handle()
	assert.Equal(t, h1, -2)
	r = <-pool.CloseLedger(h1)
	assert.NoError(t, r.Err())

	r = <-pool.OpenLedger("FINDY_MEM_LEDGER")
	assert.NoError(t, r.Err())
	h2 := r.Handle()
	assert.Equal(t, h2, -1)
}

func TestListPlugins(t *testing.T) {
	names := pool.ListPlugins()
	assert.Len(t, names, 2)
}

func TestList(t *testing.T) {
	r := <-pool.List()
	assert.NoError(t, r.Err())
}
