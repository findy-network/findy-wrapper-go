package anoncreds

import (
	"flag"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/lainio/err2"
	"github.com/stretchr/testify/assert"

	"github.com/findy-network/findy-wrapper-go"
	"github.com/findy-network/findy-wrapper-go/did"
	"github.com/findy-network/findy-wrapper-go/helpers"
	"github.com/findy-network/findy-wrapper-go/ledger"
)

func TestIssuerCreateSchema(t *testing.T) {
	err2.Check(flag.Set("logtostderr", "true"))
	err2.Check(flag.Set("v", "3"))

	ut := time.Now().Unix() - 1558884840
	schemaName := fmt.Sprintf("NEW_SCHEMA_%v", ut)

	pool := helpers.OpenTestPool(t)
	w1, name1 := helpers.CreateAndOpenTestWallet(t)
	// Create DIDs
	r := <-did.CreateAndStore(w1, did.Did{Seed: "000000000000000000000000Steward1"})
	if nil != r.Err() {
		t.Errorf("did.CreateAndStore: create steward() = %v", r.Error())
	}
	stewardDID := r.Str1()
	w1Key := r.Str2()

	err := ledger.WriteDID(pool, w1, stewardDID, stewardDID, w1Key, findy.NullString,
		findy.NullString)

	w2, name2 := helpers.CreateAndOpenTestWallet(t)
	r = <-did.Create(w2)
	if nil != r.Err() {
		t.Errorf("did.Create: create steward() = %v", r.Error())
	}
	w2DID := r.Str1()
	w2Key := r.Str2()

	// try to write prover's DID to ledger even it's not need and that's why
	// we don't care the error status
	_ = ledger.WriteDID(pool, w1, stewardDID, w2DID, w2Key, findy.NullString,
		findy.NullString)

	type args struct {
		did       string
		name      string
		version   string
		attrNames string
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"1st", args{stewardDID, schemaName, "1.0", "[\"email\"]"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ==============================================================
			// Build and publish CRED SCHEMA, could be issuer or who ever
			r := <-IssuerCreateSchema(tt.args.did, tt.args.name, tt.args.version, tt.args.attrNames)
			assert.NoError(t, r.Err())
			sid := r.Str1()
			//fmt.Println(sid)
			scJSON := r.Str2()
			//fmt.Println(scJSON)
			err = ledger.WriteSchema(pool, w1, stewardDID, scJSON)
			assert.NoError(t, err)

			time.Sleep(1 * time.Second) // let ledger build everything ready

			// Read SCHEMA from Ledger
			sid, scJSON, err = ledger.ReadSchema(pool, stewardDID, sid)
			assert.NoError(t, err)

			// ===========================================================
			// === start from the issuer side of the table
			// ===========================================================

			// -----------------------------------------------------
			// Build CRED DEF from the schema: CredDef n : 1 Schema
			// we can reuse schemas in cred defs, use tag for naming
			r = <-IssuerCreateAndStoreCredentialDef(w1, stewardDID, scJSON,
				"MY_FIRMS_CRED_DEF", findy.NullString, findy.NullString)
			assert.NoError(t, r.Err())
			// note! that in normal PROTOCOL this should be read from ledger on
			// the prover side, but now we just transfer it to there with this
			// variable.
			cdid := r.Str1()

			cd := r.Str2()
			//fmt.Println(cd)
			// BUILD CRED_DEF OFFER
			r = <-IssuerCreateCredentialOffer(w1, cdid)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IssuerCreateCredentialOffer() = %v, want %v", got, tt.want)
			}
			credOffer := r.Str1()

			// Write CRED DEF to ledger = todo should be after creation =====
			err = ledger.WriteCredDef(pool, w1, stewardDID, cd)
			assert.NoError(t, err)

			// =================================================================
			// === switch to other side of the table, credential receiver ======
			// =================================================================

			// Create master secret to wallet
			msid := "findy_master_secret"
			r = <-ProverCreateMasterSecret(w2, msid)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProverCreateMasterSecret() = %v, want %v", got, tt.want)
			}
			msid = r.Str1()
			//fmt.Println(msid)

			time.Sleep(1 * time.Second) // let ledger build everything ready

			// Get CRED DEF from the ledger
			credDefID, credDef, err := ledger.ReadCredDef(pool, w2DID, cdid)
			assert.NoError(t, err)

			// build credential request to send to back to issuer
			r = <-ProverCreateCredentialReq(w2, w2DID, credOffer, credDef, msid)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProverCreateCredentialReq() = %v, want %v", got, tt.want)
			}
			crDefReq := r.Str1() // --> would send to to other side
			crDefReqMeta := r.Str2()

			// Send cred_req to other side as response to offer -->
			// =================================================================
			// === switch to other side of the table, credential issuer   ======
			// =================================================================

			type EmailCred struct {
				Email CredDefAttr `json:"email"`
			}
			emailCred := EmailCred{}
			emailToVerify := "some.dude@findy.net"
			emailCred.Email.SetRaw(emailToVerify)
			values := dto.ToJSON(emailCred)

			r = <-IssuerCreateCredential(w1, credOffer, crDefReq, values, findy.NullString, findy.NullHandle)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IssuerCreateCredential() = %v, want %v", got, tt.want)
			}
			cred := r.Str1()

			// =================================================================
			// ==> to other side, credential receiver/holder/prover
			// =================================================================

			r = <-ProverStoreCredential(w2, findy.NullString, crDefReqMeta, cred, credDef, findy.NullString)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProverStoreCredential() = %v, want %v", got, tt.want)
			}

			// =================================================================
			// == START present proof PROTOCOL here
			// == this starts EXCEPTIONALLY at the prover's side, in this
			// == test case
			// =================================================================

			attrInfo := AttrInfo{
				Name: "email",
			}
			reqAttrs := map[string]AttrInfo{
				"attr1_referent": attrInfo,
			}
			pReq := ProofRequest{
				Name:                "FirstProofReq",
				Version:             "0.1",
				Nonce:               "12345678901234567890",
				RequestedAttributes: reqAttrs,
				RequestedPredicates: map[string]PredicateInfo{},
			}
			pReqStr := dto.ToJSON(pReq)
			//fmt.Println(pReqStr)
			r = <-ProverSearchCredentialsForProofReq(w2, pReqStr, findy.NullString)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProverSearchCredentialsForProofReq() = %v, want %v", got, tt.want)
			}
			searchHandle := r.Handle()
			const fetchMax = 2
			r = <-ProverFetchCredentialsForProofReq(searchHandle, "attr1_referent", fetchMax)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProverFetchCredentialsForProofReq() = %v, want %v", got, tt.want)
			}
			credentials := r.Str1()
			//fmt.Println(credentials)
			//fmt.Println("=====================")
			// Needs to be slice, len() tells how much we did read
			credInfo := make([]Credentials, fetchMax)
			dto.FromJSONStr(credentials, &credInfo)
			r = <-ProverCloseCredentialsSearchForProofReq(searchHandle)
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProverCloseCredentialsSearchForProofReq() = %v, want %v", got, tt.want)
			}
			schemaObject := map[string]interface{}{}
			dto.FromJSONStr(scJSON, &schemaObject)
			schemas := map[string]map[string]interface{}{
				sid: schemaObject,
			}
			schemasJSON := dto.ToJSON(schemas)
			//fmt.Println(schemasJSON)

			credDefObject := map[string]interface{}{}
			dto.FromJSONStr(credDef, &credDefObject)
			credDefs := map[string]map[string]interface{}{
				credDefID: credDefObject,
			}
			credDefsJSON := dto.ToJSON(credDefs)

			reqCred := RequestedCredentials{
				SelfAttestedAttributes: map[string]string{},
				RequestedAttributes: map[string]RequestedAttrObject{
					"attr1_referent": {
						CredID:    credInfo[0].CredInfo.Referent,
						Revealed:  true,
						Timestamp: nil,
					},
				},
				RequestedPredicates: map[string]RequestedPredObject{},
			}
			reqCredJSON := dto.ToJSON(reqCred)
			//fmt.Println(reqCredJSON)
			r = <-ProverCreateProof(w2, pReqStr, reqCredJSON, msid, schemasJSON, credDefsJSON, "{}")
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProverCreateProof() = %v, want %v", got, tt.want)
			}
			proofJS := r.Str1()
			//fmt.Println(proofJS)

			r = <-VerifierVerifyProof(pReqStr, proofJS, schemasJSON, credDefsJSON, "{}", "{}")
			if got := r.Err(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifierVerifyProof() = %v, want %v", got, tt.want)
			}
			var proof Proof
			dto.FromJSONStr(proofJS, &proof)
			if !r.Yes() ||
				proof.RequestedProof.RevealedAttrs["attr1_referent"].Raw != emailToVerify {
				t.Errorf("cannot proof!")
			}
		})
	}
	helpers.CloseAndDeleteTestWallet(w2, name2, t)
	helpers.CloseAndDeleteTestWallet(w1, name1, t)
	helpers.CloseTestPool(pool, t)
}

func TestCredDefAttr_SetRawAries(t *testing.T) {
	type fields struct {
		Raw     string
		Encoded string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"32-bit", fields{"12345", "12345"}, args{"12345"}, "12345"},
		{"32-bit max", fields{"4294967294", "4294967294"}, args{"4294967294"}, "4294967294"},
		{"32-bit + 1", fields{"4294967295", "246489516295375291562292"}, args{"4294967295"}, "112186183305887719174145980249606273079028449613853036463221369477890411334752"},
		{"String", fields{"AliceGarcia", "79092137925437828096026976"}, args{"AliceGarcia"}, "49910393468956597094487206590662454017394830827292221297761729206446655031806"},
		{"Ssn Str", fields{"123-45-6789", "59474427526494974098487352"}, args{"123-45-6789"}, "744326867119662813058574151710572260086480987778735990385444735594385781152"},
		{"ACA-py sample:alice smith", fields{"", ""}, args{"Alice Smith"}, "62816810226936654797779705000772968058283780124309077049681734835796332704413"},
		{"ACA-py sample 2", fields{"", ""}, args{"2018-05-28"}, "23402637423876324098256519317695433196813217785795317220680415812348801086586"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &CredDefAttr{
				Raw:     tt.fields.Raw,
				Encoded: tt.fields.Encoded,
			}
			if got := a.SetRawAries(tt.args.s); got != tt.want {
				t.Errorf("CredDefAttr.SetRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredDefAttr_SetRaw(t *testing.T) {
	type fields struct {
		Raw     string
		Encoded string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"32-bit", fields{"12345", "12345"}, args{"12345"}, "12345"},
		{"32-bit max", fields{"4294967294", "4294967294"}, args{"4294967294"}, "4294967294"},
		{"32-bit + 1", fields{"4294967295", "246489516295375291562292"}, args{"4294967295"}, "246489516295375291562292"},
		{"String", fields{"AliceGarcia", "79092137925437828096026976"}, args{"AliceGarcia"}, "79092137925437828096026976"},
		{"Ssn Str", fields{"123-45-6789", "59474427526494974098487352"}, args{"123-45-6789"}, "59474427526494974098487352"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &CredDefAttr{
				Raw:     tt.fields.Raw,
				Encoded: tt.fields.Encoded,
			}
			if got := a.SetRaw(tt.args.s); got != tt.want {
				t.Errorf("CredDefAttr.SetRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredDefAttr_SetEncoded(t *testing.T) {
	type fields struct {
		Raw     string
		Encoded string
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"32-bit", fields{"12345", "12345"}, args{"12345"}, "12345"},
		{"32-bit max", fields{"4294967294", "4294967294"}, args{"4294967294"}, "4294967294"},
		{"32-bit + 1", fields{"4294967295", "246489516295375291562292"}, args{"246489516295375291562292"}, "4294967295"},
		{"String", fields{"AliceGarcia", "79092137925437828096026976"}, args{"79092137925437828096026976"}, "AliceGarcia"},
		{"Ssn Str", fields{"123-45-6789", "59474427526494974098487352"}, args{"59474427526494974098487352"}, "123-45-6789"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &CredDefAttr{
				Raw:     tt.fields.Raw,
				Encoded: tt.fields.Encoded,
			}
			if got := a.SetEncoded(tt.args.s); got != tt.want {
				t.Errorf("CredDefAttr.SetEncoded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProofReq(t *testing.T) {
	proofReqJSON := `
{
  "name": "Proof of Education",
  "version": "1.0",
  "nonce": "21086678806226413236255057234572685337",
  "requested_attributes": {
    "0_name_uuid": {
      "name": "name",
      "restrictions": [
        {
          "issuer_did": "5Lgx4KLRTNgexDqT7WDALu"
        }
      ]
    },
    "0_date_uuid": {
      "name": "date",
      "restrictions": [
        {
          "issuer_did": "5Lgx4KLRTNgexDqT7WDALu"
        }
      ]
    },
    "0_degree_uuid": {
      "name": "degree",
      "restrictions": [
        {
          "issuer_did": "5Lgx4KLRTNgexDqT7WDALu"
        }
      ]
    },
    "0_self_attested_thing_uuid": {
      "name": "self_attested_thing"
    }
  },
  "requested_predicates": {
    "0_age_GE_uuid": {
      "name": "age",
      "p_type": ">=",
      "p_value": 18,
      "restrictions": [
        {
          "issuer_did": "5Lgx4KLRTNgexDqT7WDALu"
        }
      ]
    }
  }
}
`

	var proofReq ProofRequest
	dto.FromJSONStr(proofReqJSON, &proofReq)

	if proofReq.Name != "Proof of Education" {
		fmt.Print("proof req:", proofReqJSON)
		t.Error("cannot read proof request")
	}
}
