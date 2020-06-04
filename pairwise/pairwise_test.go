package pairwise

import (
	"reflect"
	"testing"

	"github.com/optechlab/findy-go"
	"github.com/optechlab/findy-go/did"
	"github.com/optechlab/findy-go/dto"
	"github.com/optechlab/findy-go/helpers"
	"github.com/optechlab/findy-go/ledger"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndExists(t *testing.T) {
	p := helpers.OpenTestPool(t)
	w1, wn1 := helpers.CreateAndOpenTestWallet(t)
	w2, wn2 := helpers.CreateAndOpenTestWallet(t)

	type args struct {
		metaData string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"1st", args{metaData: "{\"test\":\"abcd\"}"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create DIDs
			r := <-did.CreateAndStore(w1, did.Did{Seed: "000000000000000000000000Steward1"})
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAndStore() = %v, want %v", got, tt.want)
			}
			w1DID := r.Str1()
			w1Key := r.Str2()
			// Set endpoint address for 1st DID
			if got := (<-did.SetEndpoint(w1, w1DID, "this.is.test.address", w1Key)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetEndpoint() = %v, want %v", got, tt.want)
			}
			r = <-did.Endpoint(w1, p, w1DID)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint() = %v, want %v", got, tt.want)
			}
			//fmt.Println(r.Str1(), r.Str2(), r.Str3())

			// commit 1st DID to ledger to publish it's data
			err := ledger.WriteDID(p, w1, w1DID, w1DID, w1Key,
				findy.NullString, findy.NullString)
			assert.NoError(t, err)

			// create 2nd DID to wallet 2
			r = <-did.CreateAndStore(w2, did.Did{Seed: ""})
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAndStore() = %v, want %v", got, tt.want)
			}
			w2DID := r.Str1()
			w2Key := r.Str2()

			// Store THEIR DIDs to "my" wallet
			r = <-did.StoreTheir(w1, dto.ToJSON(did.Did{Did: w2DID, VerKey: w2Key}))
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StoreTheir() = %v, want %v", got, tt.want)
			}
			r = <-did.StoreTheir(w2, dto.ToJSON(did.Did{Did: w1DID, VerKey: w1Key}))
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StoreTheir() = %v, want %v", got, tt.want)
			}

			// =======================================
			// Ledger wont transfer endpoints but we have to set them by our selves during the pairwise operation
			if got := (<-did.SetEndpoint(w2, w1DID, "this.is.test.address", w1Key)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetEndpoint() = %v, want %v", got, tt.want)
			}
			// See what is in Their DIDs endpoint
			r = <-did.Endpoint(w2, p, w1DID)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint() = %v, want %v", got, tt.want)
			}
			//fmt.Println("THEIR ENDPOINT DATA:")
			//fmt.Println(r.Str1(), r.Str2(), r.Str3())
			// =======================================

			// Create pairwise for DIDs for both Wallets
			if got := (<-Create(w1, w2DID, w1DID, tt.args.metaData)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
			if got := (<-Create(w2, w1DID, w2DID, tt.args.metaData)).Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}

			// Check the pairwise existence
			if !(<-Exists(w1, w2DID)).Yes() {
				t.Errorf("Exist2")
			}
			if !(<-Exists(w2, w1DID)).Yes() {
				t.Errorf("Exist2")
			}

			// Get and print pairwise data for both "their" DIDs
			r = <-Get(w1, w2DID)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
			//d1 := r.Str1()
			//fmt.Println(d1)
			r = <-Get(w2, w1DID)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
			//l1 := r.Str1()
			//fmt.Println(l1)

			// Set meta data for "their" DID in both wallets
			r = <-did.SetMeta(w1, w2DID, dto.ToJSON(did.Did{Did: w2DID, VerKey: w2Key}))
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("did.SetMeta() = %v, want %v", got, tt.want)
			}
			r = <-did.SetMeta(w2, w1DID, dto.ToJSON(did.Did{Did: w1DID, VerKey: w1Key}))
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("did.SetMeta() = %v, want %v", got, tt.want)
			}

			// Get and print meta data for both "their" DIDs
			r = <-did.Meta(w1, w2DID)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("did.Meta() = %v, want %v", got, tt.want)
			}
			//d1 = r.Str1()
			//l1 = r.Str2()
			//fmt.Println(d1)
			//fmt.Println(l1)
			r = <-did.Meta(w2, w1DID)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("did.Meta() = %v, want %v", got, tt.want)
			}
			//l1 = r.Str1()
			//fmt.Println(l1)
			//l1 = r.Str2()
			//fmt.Println(l1)

			// All pairwise data for second wallet
			r = <-List(w2)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateAndStore() = %v, want %v", got, tt.want)
			}
			//l2 := r.Str1()
			//fmt.Println(l2)
		})
	}
	helpers.CloseAndDeleteTestWallet(w1, wn1, t)
	helpers.CloseAndDeleteTestWallet(w2, wn2, t)
	helpers.CloseTestPool(p, t)
}
