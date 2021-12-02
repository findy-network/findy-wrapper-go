package pool_test

import (
	"reflect"
	"testing"

	_ "github.com/findy-network/findy-wrapper-go/addons"
	"github.com/findy-network/findy-wrapper-go/pool"
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
	r := <-pool.OpenLedger("FINDY_MEM_LEDGER", "", "FINDY_ECHO_LEDGER", "")
	assert.NoError(t, r.Err())
	h1 := r.Handle()
	assert.Equal(t, h1, -2)

	r = <-pool.CloseLedger(h1)
	assert.NoError(t, r.Err())
}

func TestCloseLedger(t *testing.T) {
	r := <-pool.OpenLedger("FINDY_MEM_LEDGER", "", "FINDY_ECHO_LEDGER", "")
	assert.NoError(t, r.Err())
	h1 := r.Handle()
	assert.Equal(t, h1, -2)
	r = <-pool.CloseLedger(h1)
	assert.NoError(t, r.Err())

	r = <-pool.OpenLedger("FINDY_MEM_LEDGER", "")
	assert.NoError(t, r.Err())
	h2 := r.Handle()
	assert.Equal(t, h2, -1)
}

func TestList(t *testing.T) {
	r := <-pool.List()
	assert.NoError(t, r.Err())
}

func TestConverPluginArgs(t *testing.T) {
	tests := []struct {
		name   string
		arg    string
		result []string
	}{
		{"only real ledger name",
			"von",
			[]string{"von"}},
		{"plugin and name",
			"FINDY_LEDGER,von",
			[]string{"FINDY_LEDGER", "von"}},
		{"plugin and name",
			"FINDY_LEDGER,von,FINDY_MEM_LEDGER,cache",
			[]string{"FINDY_LEDGER", "von", "FINDY_MEM_LEDGER", "cache"},
		},
	}
	for _, tt := range tests {
		pools := pool.ConvertPluginArgs(tt.arg)
		assert.Equal(t, tt.result, pools)
	}
}

func TestBuildLegacyPluginArgs(t *testing.T) {
	tests := []struct {
		name   string
		arg    []string
		result []string
	}{
		{"only real ledger name",
			[]string{"test"},
			[]string{"FINDY_LEDGER", "test"},
		},
		{"only real ledger name",
			[]string{"von"},
			[]string{"FINDY_LEDGER", "von"},
		},
		{"plugin and name",
			[]string{"FINDY_LEDGER", "von"},
			[]string{"FINDY_LEDGER", "von"}},
		{"plugin and name",
			[]string{"FINDY_LEDGER", "von", "FINDY_MEM_LEDGER", "cache"},
			[]string{"FINDY_LEDGER", "von", "FINDY_MEM_LEDGER", "cache"},
		},
	}
	for _, tt := range tests {
		pools := pool.BuildLegacyPluginArgs(tt.arg)
		assert.Equal(t, tt.result, pools)
	}
}
