package crypto

import (
	"reflect"
	"testing"

	"github.com/findy-network/findy-wrapper-go"
	"github.com/findy-network/findy-wrapper-go/did"
	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/findy-network/findy-wrapper-go/helpers"
)

var connectionRequest = `  {
    "@type": "did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/connections/1.0/request",
    "@id": "670bc804-2c06-453c-aee6-48d3c929b488",
    "label": "Alice Agent",
    "connection": {
      "DID": "ERYihzndieTdh4UA7Q6Y3C",
      "DIDDoc": {
        "@context": "https://w3id.org/did/v1",
        "id": "did:sov:ERYihzndieTdh4UA7Q6Y3C",
        "publicKey": [
          {
            "id": "did:sov:ERYihzndieTdh4UA7Q6Y3C#1",
            "type": "Ed25519VerificationKey2018",
            "controller": "did:sov:ERYihzndieTdh4UA7Q6Y3C",
            "publicKeyBase58": "8KLQJNs7cJFY5vcRTWzb33zYr5zhDrcaX6jgD5Uaofcu"
          }
        ],
        "authentication": [
          {
            "type": "Ed25519SignatureAuthentication2018",
            "publicKey": "did:sov:ERYihzndieTdh4UA7Q6Y3C#1"
          }
        ],
        "service": [
          {
            "id": "did:sov:ERYihzndieTdh4UA7Q6Y3C;indy",
            "type": "IndyAgent",
            "priority": 0,
            "recipientKeys": ["8KLQJNs7cJFY5vcRTWzb33zYr5zhDrcaX6jgD5Uaofcu"],
            "serviceEndpoint": "http://192.168.65.3:8030"
          }
        ]
      }
    }
  }
`

func TestVerifySignature(t *testing.T) {
	p := helpers.OpenTestPool(t)
	w, wn := helpers.CreateAndOpenTestWallet(t)
	r := <-did.Create(w)
	if nil != r.Err() {
		t.Errorf("did.Create: create steward() = %v", r.Error())
	}
	//w1DID := r.Str1()
	w1Key := r.Str2()
	m := []byte("test message 10010")
	r = <-SignMsg(w, w1Key, m)
	if nil != r.Err() {
		t.Errorf("cannot sign message(%v) with key=%v", m, w1Key)
	}
	sig := r.Bytes()

	// Testing that there really isn't any implicit dependency to wallet, only
	// for verkey, so, cleaning up already.
	helpers.CloseAndDeleteTestWallet(w, wn, t)
	helpers.CloseTestPool(p, t)

	r = <-VerifySignature(w1Key, m, sig)
	if nil != r.Err() || !r.Yes() {
		t.Errorf("cannot verify signature(%v) with key=%v", sig, w1Key)
	}
}

func TestPackMessage(t *testing.T) {
	//ctx.SetTrace(true)

	p := helpers.OpenTestPool(t)
	w, wn := helpers.CreateAndOpenTestWallet(t)
	r := <-did.Create(w)
	if nil != r.Err() {
		t.Errorf("did.Create: create steward() = %v", r.Error())
	}
	//w1DID := r.Str1()
	w1Key := r.Str2()

	keys := dto.JSONArray(w1Key)
	//keys := fmt.Sprintf("[\"%s\"]", w1Key)

	w2, w2n := helpers.CreateAndOpenTestWallet(t)
	r = <-did.Create(w2)
	if nil != r.Err() {
		t.Errorf("did.Create: create steward() = %v", r.Error())
	}
	//w2DID := r.Str1()
	w2Key := r.Str2()

	type args struct {
		wallet            int
		wallet2           int
		recipientKeysJSON string
		senderKey         string
		msg               []byte
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"1st", args{w2, w, keys, w2Key, []byte(connectionRequest)}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// No sender key, aka anonCrypt, use wrong wallet handle to see it's ignored
			r := <-Pack(-1, findy.NullString, tt.args.msg, w1Key)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pack() = %v, want %v", got, tt.want)
			}
			bytes := r.Bytes()
			r = <-UnpackMessage(tt.args.wallet2, bytes)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackMessage() = %v, want %v", got, tt.want)
			}

			// Rest with the sender key
			r = <-Pack(tt.args.wallet, tt.args.senderKey, tt.args.msg, w1Key)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pack() = %v, want %v", got, tt.want)
			}
			bytes = r.Bytes()
			r = <-UnpackMessage(tt.args.wallet2, bytes)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnpackMessage() = %v, want %v", got, tt.want)
			}
			unp := NewUnpacked(r.Bytes())
			if !reflect.DeepEqual(unp.Bytes(), tt.args.msg) {
				t.Errorf("UnpackMessage() = %v, want %v", unp.Bytes(), tt.args.msg)
			}

			r = <-PackMessage(tt.args.wallet, tt.args.recipientKeysJSON, tt.args.senderKey, tt.args.msg)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PackMessage() = %v, want %v", got, tt.want)
			}
			bytes = r.Bytes()
			up := NewUnpacking(tt.args.wallet2, bytes)
			rBytes, err := up.Bytes()
			if err != nil || !reflect.DeepEqual(rBytes, tt.args.msg) {
				t.Errorf("UnpackMessage() = %v, want %v", rBytes, tt.args.msg)
			}

		})
	}

	helpers.CloseAndDeleteTestWallet(w2, w2n, t)
	helpers.CloseAndDeleteTestWallet(w, wn, t)
	helpers.CloseTestPool(p, t)
}
