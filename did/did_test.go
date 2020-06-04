package did

import (
	"fmt"

	"reflect"
	"testing"

	"github.com/findy-network/findy-wrapper-go/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndStore(t *testing.T) {
	//p := helpers.OpenTestPool(t)
	w, wn := helpers.CreateAndOpenTestWallet(t)

	type args struct {
		wallet int
		did    Did
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"1st", args{wallet: w, did: Did{Seed: ""}}, nil},
		{"2nd", args{wallet: w, did: Did{Seed: ""}}, nil},
		{"3nd", args{wallet: w, did: Did{Seed: ""}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := <-CreateAndStore(tt.args.wallet, tt.args.did)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAndStore() = %v, want %v", got, tt.want)
			}
			didStr := r.Str1()
			didKey := r.Str2()
			if got := (<-SetEndpoint(w, didStr, "this.is.test.address", didKey)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetEndpoint() = %v, want %v", got, tt.want)
			}
			//r = <-Endpoint(w, p, didStr)
			r = <-Endpoint(w, -1, didStr)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint() = %v, want %v", got, tt.want)
			}
			fmt.Println(r.Str1(), r.Str2(), r.Str3())
		})
	}
	helpers.CloseAndDeleteTestWallet(w, wn, t)
	//helpers.CloseTestPool(p, t)
}

func TestKey(t *testing.T) {
	w, wn := helpers.CreateAndOpenTestWallet(t)

	r := <-CreateAndStore(w, Did{Seed: ""})
	assert.NoError(t, r.Err())
	did := r.Str1()
	vk := r.Str2()

	r = <-Key(-1, w, did)
	assert.NoError(t, r.Err())
	assert.Equal(t, vk, r.Str1())

	helpers.CloseAndDeleteTestWallet(w, wn, t)
}
