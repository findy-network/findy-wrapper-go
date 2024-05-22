package anoncreds

import (
	"fmt"
	"testing"
	"time"

	"github.com/findy-network/findy-wrapper-go/dto"
	"github.com/golang/glog"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"

	"github.com/findy-network/findy-wrapper-go"
	"github.com/findy-network/findy-wrapper-go/did"
	"github.com/findy-network/findy-wrapper-go/helpers"
	"github.com/findy-network/findy-wrapper-go/ledger"
)

const ledgerWaitTimer = 1 * time.Second

func TestIssuerCreateSchema(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
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
	assert.NoError(err)

	err = ledger.WriteDID(pool, w1, stewardDID, stewardDID, w1Key, findy.NullString,
		findy.NullString)
	assert.NoError(err, "updates are fine for now")

	w2, name2 := helpers.CreateAndOpenTestWallet(t)
	r = <-did.Create(w2)
	if nil != r.Err() {
		t.Errorf("did.Create: create steward() = %v", r.Error())
	}
	w2DID := r.Str1()
	w2Key := r.Str2()

	err = ledger.WriteDID(pool, w1, stewardDID, w2DID, w2Key, findy.NullString,
		findy.NullString)
	// try to write prover's DID to ledger even it's not need and that's why
	// we KNOW that there shouldn'd be any errors.
	assert.NoError(err)

	type args struct {
		did       string
		name      string
		version   string
		attrNames string
	}
	tests := []struct {
		name  string
		args  args
		noErr bool
	}{
		{"1st", args{stewardDID, schemaName, "1.0", "[\"email\"]"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer assert.PushTester(t)()

			// ==============================================================
			// Holder side:
			// Create master secret to wallet for Holder wallet
			msid := "findy_master_secret"
			r = <-ProverCreateMasterSecret(w2, msid)
			try.To(r.Err())
			msid = r.Str1()
			// ==============================================================

			// ==============================================================
			// Issuer side:
			// Build and publish CRED SCHEMA, could be issuer or who ever
			r := <-IssuerCreateSchema(tt.args.did, tt.args.name, tt.args.version, tt.args.attrNames)
			assert.NoError(r.Err())
			sid := r.Str1()
			scJSON := r.Str2()

			err = ledger.WriteSchema(pool, w1, stewardDID, scJSON)
			assert.NoError(err)

			glog.V(2).Infoln("<<====IN SchemaID:", sid, "waiting before read schema")
			time.Sleep(ledgerWaitTimer) // let ledger build everything ready

			// Read SCHEMA from Ledger
			sid, scJSON, err = ledger.ReadSchema(pool, stewardDID, sid)
			assert.NoError(err)
			glog.V(2).Infoln("=====================> OUT SchemaID:", sid)
			glog.V(2).Infoln("==== getting from ledger for schema:", scJSON)

			// -----------------------------------------------------
			// Build CRED DEF from the schema: CredDef n : 1 Schema
			// we can reuse schemas in cred defs, use tag for naming
			r = <-IssuerCreateAndStoreCredentialDef(w1, stewardDID, scJSON,
				"MY_FIRMS_CRED_DEF", findy.NullString, findy.NullString)
			assert.NoError(r.Err())
			// note! that in normal PROTOCOL this should be read from ledger on
			// the prover side, but now we just transfer it to there with this
			// variable.
			cdid := r.Str1()
			cd := r.Str2()
			// Write CRED DEF to ledger
			glog.V(2).Infoln("<<==================== IN CredDefID:", cdid)
			try.To(ledger.WriteCredDef(pool, w1, stewardDID, cd))

			// ==============================================================
			// == BUILD CRED_DEF OFFER
			// ==============================================================
			r = <-IssuerCreateCredentialOffer(w1, cdid)
			try.To(r.Err())
			credOffer := r.Str1()

			// =================================================================
			// === switch to other side of the table, credential receiver ======
			// =================================================================

			time.Sleep(ledgerWaitTimer) // let ledger build everything ready

			// Get CRED DEF from the ledger
			credDefID, credDef, err := ledger.ReadCredDef(pool, w2DID, cdid)
			assert.NoError(err)
			glog.V(2).Infoln("<<==================== OUT CredDef:", credDefID)

			// Send cred_req to other side as response to offer -->
			// =================================================================
			// === switch to other side of the table, credential issuer   ======
			// =================================================================

			var emailToVerify string
			for i := 0; i < 2; i++ {
				// build credential REQUEST to send to back to issuer
				r = <-ProverCreateCredentialReq(w2, w2DID, credOffer, credDef, msid)
				try.To(r.Err())
				crDefReq := r.Str1() // --> would send to to other side
				crDefReqMeta := r.Str2()

				type EmailCred struct {
					Email CredDefAttr `json:"email"`
				}
				emailCred := EmailCred{}
				emailToVerify = fmt.Sprintf("some%d.dude@findy.net", i)
				emailCred.Email.SetRaw(emailToVerify)
				values := dto.ToJSON(emailCred)

				r = <-IssuerCreateCredential(w1, credOffer, crDefReq, values, findy.NullString, findy.NullHandle)
				try.To(r.Err())
				cred := r.Str1()

				// =================================================================
				// ==> to other side, credential receiver/holder/prover
				// =================================================================

				r = <-ProverStoreCredential(w2, findy.NullString, crDefReqMeta, cred, credDef, findy.NullString)
				try.To(r.Err())
				glog.V(2).Infoln("======== Credential: ", i, " ", emailToVerify)
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
			wql := fmt.Sprintf(`"attr::email::value": "%s"`, emailToVerify)
			wqlJSONStr := fmt.Sprintf(`{"attr1_referent": {%s}}`, wql)
			glog.V(2).Infoln("---------")
			glog.V(2).Infoln("--", wqlJSONStr)
			glog.V(2).Infoln("---------")
			r = <-ProverSearchCredentialsForProofReq(w2, pReqStr, wqlJSONStr)
			try.To(r.Err())
			searchHandle := r.Handle()
			const fetchMax = 4
			r = <-ProverFetchCredentialsForProofReq(searchHandle, "attr1_referent", fetchMax)
			try.To(r.Err())
			credentials := r.Str1()
			glog.V(2).Infoln(credentials)
			// Needs to be slice, len() tells how much we did read
			credInfo := make([]Credentials, 0, fetchMax)
			dto.FromJSONStr(credentials, &credInfo)
			assert.SLen(credInfo, 1)
			glog.V(2).Infoln("--> len: ", len(credInfo))
			r = <-ProverCloseCredentialsSearchForProofReq(searchHandle)
			try.To(r.Err())
			schemaObject := map[string]interface{}{}
			dto.FromJSONStr(scJSON, &schemaObject)
			schemas := map[string]map[string]interface{}{
				sid: schemaObject,
			}
			schemasJSON := dto.ToJSON(schemas)

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
			r = <-ProverCreateProof(w2, pReqStr, reqCredJSON, msid, schemasJSON, credDefsJSON, "{}")
			try.To(r.Err())
			proofJS := r.Str1()

			r = <-VerifierVerifyProof(pReqStr, proofJS, schemasJSON, credDefsJSON, "{}", "{}")
			try.To(r.Err())
			var proof Proof
			dto.FromJSONStr(proofJS, &proof)
			assert.That(r.Yes())
			assert.Equal(proof.RequestedProof.RevealedAttrs["attr1_referent"].Raw, emailToVerify)
		})
	}
	helpers.CloseAndDeleteTestWallet(w2, name2, t)
	helpers.CloseAndDeleteTestWallet(w1, name1, t)
	helpers.CloseTestPool(pool, t)
}

func TestCredDefAttr_SetRawAries(t *testing.T) {
	assert.PushTester(t)
	defer assert.PopTester()
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
			assert.PushTester(t)
			defer assert.PopTester()
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
	assert.PushTester(t)
	defer assert.PopTester()
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
			assert.PushTester(t)
			defer assert.PopTester()
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
	assert.PushTester(t)
	defer assert.PopTester()
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
			assert.PushTester(t)
			defer assert.PopTester()
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
	assert.PushTester(t)
	defer assert.PopTester()
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
		t.Error("cannot read proof request")
	}
}
