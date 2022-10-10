/*
Package anoncreds is corresponding Go package for libindy's anoncreds namespace.
We suggest that you read indy SDK documentation for more information.
Unfortunately the documentation during the writing of this package was minimal
to nothing. We suggest you do the same as we did, read the rust code if more
detailed information is needed.
*/
package anoncreds

import "C"
import (
	"crypto/sha256"
	"math"
	"math/big"

	"github.com/findy-network/findy-wrapper-go/internal/c2go"
	"github.com/findy-network/findy-wrapper-go/internal/ctx"
)

// These are for indy/anoncreds.rs which seems to be the only specs for these
// JSON files.
//
/// filter_json: filter for credentials
///        {
///            "schema_id": string, (Optional)
///            "schema_issuer_did": string, (Optional)
///            "schema_name": string, (Optional)
///            "schema_version": string, (Optional)
///            "issuer_did": string, (Optional)
///            "cred_def_id": string, (Optional)
///        }
/// cb: Callback that takes command result as parameter.
///
/// #Returns
/// credentials json
///     [{
///         "referent": string, // cred_id in the wallet
///         "attrs": {"key1":"raw_value1", "key2":"raw_value2"},
///         "schema_id": string,
///         "cred_def_id": string,
///         "rev_reg_id": Optional<string>,
///         "cred_rev_id": Optional<string>
///     }]
///

/// proof_request_json: proof request json
///     {
///         "name": string,
///         "version": string,
///         "nonce": string,
///         "requested_attributes": { // set of requested attributes
///              "<attr_referent>": <attr_info>, // see below
///              ...,
///         },
///         "requested_predicates": { // set of requested predicates
///              "<predicate_referent>": <predicate_info>, // see below
///              ...,
///          },
///         "non_revoked": Optional<<non_revoc_interval>>, // see below,
///                        // If specified prover must proof non-revocation
///                        // for date in this interval for each attribute
///                        // (can be overridden on attribute level)
///     }
/// cb: Callback that takes command result as parameter.
///
/// where
/// attr_referent: Proof-request local identifier of requested attribute
/// attr_info: Describes requested attribute
///     {
///         "name": string, // attribute name, (case insensitive and ignore spaces)
///         "restrictions": Optional<filter_json>, // see above
///         "non_revoked": Optional<<non_revoc_interval>>, // see below,
///                        // If specified prover must proof non-revocation
///                        // for date in this interval this attribute
///                        // (overrides proof level interval)
///     }
/// predicate_referent: Proof-request local identifier of requested attribute predicate
/// predicate_info: Describes requested attribute predicate
///     {
///         "name": attribute name, (case insensitive and ignore spaces)
///         "p_type": predicate type (Currently ">=" only)
///         "p_value": int predicate value
///         "restrictions": Optional<filter_json>, // see above
///         "non_revoked": Optional<<non_revoc_interval>>, // see below,
///                        // If specified prover must proof non-revocation
///                        // for date in this interval this attribute
///                        // (overrides proof level interval)
///     }
/// non_revoc_interval: Defines non-revocation interval
///     {
///         "from": Optional<int>, // timestamp of interval beginning
///         "to": Optional<int>, // timestamp of interval ending
///     }
///
/// #Returns
/// credentials_json: json with credentials for the given proof request.
///     {
///         "requested_attrs": {
///             "<attr_referent>": [{ cred_info: <credential_info>, interval: Optional<non_revoc_interval> }],
///             ...,
///         },
///         "requested_predicates": {
///             "requested_predicates": [{ cred_info: <credential_info>, timestamp: Optional<integer> },
///                 { cred_info: <credential_2_info>, timestamp: Optional<integer> }],
///             "requested_predicate_2_referent": [{ cred_info: <credential_2_info>, timestamp: Optional<integer> }]
///         }
///     }, where credential is
///     {
///         "referent": <string>,
///         "attrs": {"attr_name" : "attr_raw_value"},
///         "schema_id": string,
///         "cred_def_id": string,
///         "rev_reg_id": Optional<int>,
///         "cred_rev_id": Optional<int>,
///     }
///

// Proof is wrapper struct for libindy's corresponding JSON type.
type Proof struct {
	RequestedProof RequestedProof         `json:"requested_proof"`
	Proof          map[string]interface{} `json:"proof"`
	Identifiers    []IdentifiersObj       `json:"identifiers"`
}

// IdentifiersObj is wrapper struct for libindy's corresponding JSON type.
type IdentifiersObj struct {
	SchemaID  string `json:"schema_id"`
	CredDefID string `json:"cred_def_id"`
	RevRegID  string `json:"rev_reg_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// RequestedProof is wrapper struct for libindy's corresponding JSON type.
type RequestedProof struct {
	RevealedAttrs     map[string]RevealedAttr `json:"revealed_attrs"`
	UnrevealedAttrs   map[string]interface{}  `json:"unrevealed_attrs"`
	SelfAttestedAttrs map[string]interface{}  `json:"self_attested_attrs"`
	Predicates        map[string]interface{}  `json:"predicates"`
}

// RevealedAttr is wrapper struct for libindy's corresponding JSON type.
type RevealedAttr struct {
	SubProofIndex int    `json:"sub_proof_index"`
	Raw           string `json:"raw"`
	Encoded       string `json:"encoded"`
}

// RequestedCredentials is wrapper struct for libindy's corresponding JSON type.
type RequestedCredentials struct {
	SelfAttestedAttributes map[string]string              `json:"self_attested_attributes"`
	RequestedAttributes    map[string]RequestedAttrObject `json:"requested_attributes"`
	RequestedPredicates    map[string]RequestedPredObject `json:"requested_predicates"`
}

// RequestedPredObject is wrapper struct for libindy's corresponding JSON type.
type RequestedPredObject struct {
	CredID    string `json:"cred_id"`
	Timestamp *int   `json:"timestamp,omitempty"`
}

// RequestedAttrObject is wrapper struct for libindy's corresponding JSON type.
type RequestedAttrObject struct {
	CredID    string `json:"cred_id"`
	Timestamp *int   `json:"timestamp"`
	Revealed  bool   `json:"revealed"`
}

// Credentials is wrapper struct for libindy's corresponding JSON type.
type Credentials struct {
	CredInfo CredentialInfo    `json:"cred_info"`
	Interval *NonRevocInterval `json:"interval,omitempty"`
}

// CredentialInfo is wrapper struct for libindy's corresponding JSON type.
type CredentialInfo struct {
	Referent  string            `json:"referent"`
	Attrs     map[string]string `json:"attrs"`
	SchemaID  string            `json:"schema_id"`
	CredDefID string            `json:"cred_def_id"`
	RevRegID  int               `json:"rev_reg_id,omitempty"`
	CredRevID int               `json:"cred_rev_id,omitempty"`
}

// ProofRequest is wrapper struct for libindy's corresponding JSON type.
type ProofRequest struct {
	Name                string                   `json:"name"`
	Version             string                   `json:"version"`
	Nonce               string                   `json:"nonce"`
	RequestedAttributes map[string]AttrInfo      `json:"requested_attributes"`
	RequestedPredicates map[string]PredicateInfo `json:"requested_predicates"`
	NonRevoked          *NonRevocInterval        `json:"non_revoked,omitempty"`
}

// PredicateInfo is wrapper struct for libindy's corresponding JSON type.
type PredicateInfo struct {
	Name         string           `json:"name"`
	PType        string           `json:"p_type"`
	PValue       int              `json:"p_value"`
	Restrictions []Filter         `json:"restrictions,omitempty"`
	NonRevoked   NonRevocInterval `json:"non_revoked,omitempty"`
}

// AttrInfo is wrapper struct for libindy's corresponding JSON type.
type AttrInfo struct {
	// attribute names, (case insensitive and ignore spaces)
	// indy-sdk NOTE: should either be "name" or "names", not both and not none of them.
	// Use "names" to specify several attributes that have to match a single credential.
	Name         string            `json:"name"`
	Names        []string          `json:"names,omitempty"`
	Restrictions []Filter          `json:"restrictions,omitempty"`
	NonRevoked   *NonRevocInterval `json:"non_revoked,omitempty"`
}

// Filter is wrapper struct for libindy's corresponding JSON type.
type Filter struct {
	SchemaID        string `json:"schema_id,omitempty"`
	SchemaIssuerDID string `json:"schema_issuer_did,omitempty"`
	SchemaName      string `json:"schema_name,omitempty"`
	IssuerDID       string `json:"issuer_did,omitempty"`
	CredDefID       string `json:"cred_def_id,omitempty"`
}

// NonRevocInterval is wrapper struct for libindy's corresponding JSON type.
type NonRevocInterval struct {
	From int `json:"from,omitempty"`
	To   int `json:"to,omitempty"`
}

// CredDefCfg is wrapper struct for libindy's corresponding JSON type.
type CredDefCfg struct {
	SupportRevocation bool `json:"support_revocation"`
}

// CredDefAttr is wrapper struct for libindy's corresponding JSON type.
type CredDefAttr struct {
	Raw     string `json:"raw"`
	Encoded string `json:"encoded"`
}

// SetRaw sets raw value to attribute. This is obsolete with Aries.
func (a *CredDefAttr) SetRaw(s string) string {
	a.Raw = s
	a.Encoded = encode(s)
	return a.Encoded
}

// SetRawAries sets the raw value of the attribute and writes the encoded value
// of it in same algorithm as Aries indy implementations
func (a *CredDefAttr) SetRawAries(s string) string {
	a.Raw = s
	a.Encoded = encodeAries(s)
	return a.Encoded
}

// SetEncoded sets encoded value to attribute. This is obsolete with Aries.
func (a *CredDefAttr) SetEncoded(s string) string {
	a.Encoded = s
	a.Raw = decode(s)
	return a.Raw
}

func encodeAries(s string) string {
	a := big.NewInt(math.MaxUint32)
	b := big.NewInt(0)
	_, ok := b.SetString(s, 10)
	if ok && b.Cmp(a) == -1 { // is number and b < a
		return s
	}
	sum := sha256.Sum256([]byte(s))
	b.SetBytes(sum[:])
	return b.String()
}

func encode(s string) string {
	a := big.NewInt(math.MaxUint32)
	b := big.NewInt(0)
	_, ok := b.SetString(s, 10)
	if ok && b.Cmp(a) == -1 { // is number and b < a
		return s
	}
	b.SetBytes([]byte(s))
	a.Add(a, b)
	return a.String()
}

func decode(s string) string {
	a := big.NewInt(0)
	a.SetString(s, 10)
	b := big.NewInt(math.MaxUint32)
	if a.Cmp(b) == -1 { // a < b
		return s
	}
	a.Sub(a, b)
	return string(a.Bytes())
}

// IssuerCreateSchema creates credential schema entity that describes
// credential attributes list and allows credentials interoperability. Schema
// is public and intended to be shared with all anoncreds workflow actors
// usually by publishing SCHEMA transaction to Indy distributed ledger.
//
// It is IMPORTANT for current version POST Schema in Ledger and after that GET
// it from Ledger with correct seq_no to save compatibility with Ledger. After
// that can call indy_issuer_create_and_store_credential_def to build
// corresponding Credential Definition.
func IssuerCreateSchema(did, name, version, attrNames string) ctx.Channel {
	return c2go.FindyIssuerCreateSchema(did, name, version, attrNames)
}

/*
IssuerCreateAndStoreCredentialDef creates and stores a credential definition
entity that encapsulates credentials issuer DID, credential schema, secrets used
for signing credentials and secrets used for credentials revocation.

issuer_did: a DID of the issuer signing cred_def transaction to the Ledger

Credential definition entity contains private and public parts. NOTE! Private
part will be stored in the wallet. Public part will be returned as JSON intended
to be shared with all anoncreds workflow actors usually by publishing CRED_DEF
transaction to Indy distributed ledger.

It is IMPORTANT for current version GET Schema from Ledger with correct
seq_no to save compatibility with Ledger.
*/
func IssuerCreateAndStoreCredentialDef(wallet int, did, schema, tag, sigType, config string) ctx.Channel {
	return c2go.FindyIssuerCreateAndStoreCredentialDef(wallet, did, schema, tag, sigType, config)
}

// IssuerCreateCredentialOffer creates on credential offer for part of indy
// specific issuing protocol. It's called before an issuer send its result as on
// offer credential.
func IssuerCreateCredentialOffer(wallet int, credDefID string) ctx.Channel {
	return c2go.FindyIssuerCreateCredentialOffer(wallet, credDefID)
}

// IssuerCreateCredential creates credential at issuer side. Two last arguments
// are optional. blobHandle should be -1 if not given.
// Returns:
// cred_json: Credential json containing signed credential values.
// cred_revoc_id: local id, which can be used for revocation of this credential.
// revoc_reg_delta_json: revocation registry delta.
func IssuerCreateCredential(wallet int, credOffer, credReq, credValues, revRegID string, blobHandle int) ctx.Channel {
	return c2go.FindyIssuerCreateCredential(wallet, credOffer, credReq, credValues, revRegID, blobHandle)
}

// ProverCreateMasterSecret creates a master secret with a given id and stores
// it in the wallet. The id must be unique.
func ProverCreateMasterSecret(wallet int, id string) ctx.Channel {
	return c2go.FindyProverCreateMasterSecret(wallet, id)
}

/*
ProverCreateProof creates a proof according to the given proof request either a
corresponding credential with optionally revealed attributes or self-attested
attribute must be provided for each requested attribute (see
ProverGetCredentialsForPoolReq). A proof request may request multiple
credentials from different schemas and different issuers. All required schemas,
public keys and revocation registries must be provided. The proof request also
contains nonce. The proof contains either proof or self-attested attribute value
for each requested attribute.
*/
func ProverCreateProof(wallet int, proofReq, reqCred, masterSec, schemas, credDef, revStates string) ctx.Channel {
	return c2go.FindyProverCreateProof(wallet, proofReq, reqCred, masterSec, schemas, credDef, revStates)
}

// ProverCreateCredentialReq creates a credential request for the given
// credential offer.
//
// The method creates a blinded master secret for a master secret identified by
// a provided name. The master secret identified by the name must be already
// stored in the secure wallet (see prover_create_master_secret) The blinded
// master secret is a part of the credential request.
func ProverCreateCredentialReq(wallet int, prover, credOffer, credDef, master string) ctx.Channel {
	return c2go.FindyProverCreateCredentialReq(wallet, prover, credOffer, credDef, master)
}

// ProverStoreCredential stores credentials, where parameters are as follow:
// credId: (opt, auto), ID by which credential will be stored in wallet.
// credReqMeta: json: a metadata created by ProverCreateCredentialReq.
// cred: credential json received from issuer.
// credDef: credential definition json related to <credDefId> in <cred>
// revRegDef: (opt) json related <rev_reg_def_id> in <cred json>
func ProverStoreCredential(wallet int, credID, credReqMeta, cred, credDef, revRegDef string) ctx.Channel {
	return c2go.FindyProverStoreCredential(wallet, credID, credReqMeta, cred, credDef, revRegDef)
}

// VerifierVerifyProof verifies a proof (of multiple credential). All required
// schemas, public keys and revocation registries must be provided.
func VerifierVerifyProof(proofReqJSON, proofJSON, schemasJSON, credDefsJSON, revRegDefsJSON, revRegsJSON string) ctx.Channel {
	return c2go.FindyVerifierVerifyProof(proofReqJSON, proofJSON, schemasJSON, credDefsJSON, revRegDefsJSON, revRegsJSON)
}

// ProverSearchCredentialsForProofReq searches for credentials matching the
// given proof request. Instead of immediately returning of fetched credentials
// this call returns searchHandle that can be used later to fetch records by
// small batches with ProverFetchCredentialsForProofReq.
func ProverSearchCredentialsForProofReq(wallet int, proofReqJSON, extraQueryJSON string) ctx.Channel {
	return c2go.FindyProverSearchCredentialsForProofReq(wallet, proofReqJSON, extraQueryJSON)
}

// ProverFetchCredentialsForProofReq fetches next credentials for the requested
// item using proof request search handle
//
//	searchHandle: Search handle created by ProverSearchCredentialsForProofReq
//	itemRef: Referent of attribute/predicate in the proof request
//	count: Count of credentials to fetch
func ProverFetchCredentialsForProofReq(searchHandle int, itemRef string, count int) ctx.Channel {
	return c2go.FindyProverFetchCredentialsForProofReq(searchHandle, itemRef, count)
}

// ProverCloseCredentialsSearchForProofReq closes the search identified by
// search handle.
func ProverCloseCredentialsSearchForProofReq(searchHandle int) ctx.Channel {
	return c2go.FindyProverCloseCredentialsSearchForProofReq(searchHandle)
}
