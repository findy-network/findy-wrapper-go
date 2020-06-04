/*
Package ledger is corresponding Go package for libindy's ledger namespace.
We suggest that you read indy SDK documentation for more information.
*/
package ledger

import (
	"github.com/optechlab/findy-go/internal/c2go"
	"github.com/optechlab/findy-go/internal/ctx"
)

/*
BuildNymRequest is a libindy wrapper function. See more information from indy
SDK documentation.

The function builds a ledger transaction message which will
be send to it with SignAndSubmitRequest. Note that in normal cases, you should
not use it for leger communication. You should use the special transaction
functions from this package (e.g. ReadSchema, WriteSchema,...).
*/
func BuildNymRequest(submitterDid, targetDid, verkey, alias, role string) ctx.Channel {
	return c2go.LedgerBuildNymRequest(submitterDid, targetDid, verkey, alias, role)
}

// SignAndSubmitRequest Signs and submits request message to validator pool.
//
// Adds submitter information to passed request json, signs it with submitter
// sign key (see wallet_sign), and sends signed request message to validator
// pool (see write_request).
//
// Note! You should us WriteSchema, WriteDID, ...
// instead.
func SignAndSubmitRequest(pool, wallet int, submitterDid, request string) ctx.Channel {
	return c2go.LedgerSignAndSubmitRequest(pool, wallet, submitterDid, request)
}

/*
BuildSchemaRequest is a libindy wrapper function. See more information from indy
SDK documentation.

The function builds a ledger transaction message which will be send to it with
SighAndSubmitRequest. Note that in normal cases, you should not use it for leger
communication. You should use the special transaction functions from this
package (e.g. ReadSchema, WriteSchema,...).
*/
func BuildSchemaRequest(submitterDid, data string) ctx.Channel {
	return c2go.LedgerBuildSchemaRequest(submitterDid, data)
}

/*
BuildGetNymRequest is a libindy wrapper function. See more information from indy
SDK documentation.

The function builds a ledger transaction message which will be send to it with
SubmitRequest. Note that in normal cases, you should not use it for leger
communication. You should use the special transaction functions from this
package (e.g. ReadSchema, WriteSchema,...).
*/
func BuildGetNymRequest(submitterDid, targetDid string) ctx.Channel {
	return c2go.FindyBuildGetNymRequest(submitterDid, targetDid)
}

// SubmitRequest publishes request message to validator pool (no signing, unlike
// sign_and_submit_request).
//
// The request is sent to the validator pool as is. It's assumed that it's
// already prepared.
//
// Note! You should us WriteSchema, WriteDID, ...
// instead.
func SubmitRequest(poolHandle int, request string) ctx.Channel {
	return c2go.FindySubmitRequest(poolHandle, request)
}

/*
BuildAttribRequest is a libindy wrapper function. See more information from indy
SDK documentation.

The function builds a ledger transaction message which will be send to it with
SignAndSubmitRequest. Note that in normal cases, you should not use it for leger
communication. You should use the special transaction functions from this
package (e.g. ReadSchema, WriteSchema,...).
*/
func BuildAttribRequest(sDID, tDID, hasc, raw, enc string) ctx.Channel {
	return c2go.FindyBuildAttribRequest(sDID, tDID, hasc, raw, enc)
}

/*
BuildGetAttribRequest is a libindy wrapper function. See more information from indy
SDK documentation.

The function builds a ledger transaction message which will be send to it with
SubmitRequest. Note that in normal cases, you should not use it for leger
communication. You should use the special transaction functions from this
package (e.g. ReadSchema, WriteSchema,...).
*/
func BuildGetAttribRequest(sDID, tDID, hasc, raw, enc string) ctx.Channel {
	return c2go.FindyBuildGetAttribRequest(sDID, tDID, hasc, raw, enc)
}

// BuildCredDefRequest builds an CRED_DEF request. Request to add a Credential
// Definition (in particular, public key), that Issuer creates for a particular
// Credential Schema.
//
// Note! You should us WriteSchema, WriteDID, ...
// instead.
func BuildCredDefRequest(submitter, data string) ctx.Channel {
	return c2go.FindyBuildCredDefRequest(submitter, data)
}

/*
BuildGetSchemaRequest is a libindy wrapper function. See more information from indy
SDK documentation.

The function builds a ledger transaction message which will
be send to it with SubmitRequest. Note that in normal cases, you should
not use it for leger communication. You should use the special transaction
functions from this package (e.g. ReadSchema, WriteSchema,...).
*/
func BuildGetSchemaRequest(submitter, id string) ctx.Channel {
	return c2go.FindyBuildGetSchemaRequest(submitter, id)
}

/*
ParseGetSchemaResponse is a libindy wrapper function. See more information from
indy SDK documentation.

The function parses the results from SubmitRequest. Note that in normal cases,
you should not use it for leger communication. You should use the special
transaction functions from this package (e.g. ReadSchema, ReadCredDef,...).
*/
func ParseGetSchemaResponse(response string) ctx.Channel {
	return c2go.FindyParseGetSchemaResponse(response)
}

/*
ParseGetCredDefResponse is a libindy wrapper function. See more information from
indy SDK documentation.

The function parses the results from SubmitRequest. Note that in normal cases,
you should not use it for leger communication. You should use the special
transaction functions from this package (e.g. ReadSchema, ReadCredDef,...).
*/
func ParseGetCredDefResponse(credDefResp string) ctx.Channel {
	return c2go.FindyParseGetCredDefResponse(credDefResp)
}

// BuildGetCredDefRequest builds a GET_CRED_DEF request. Request to get a
// Credential Definition (in particular, public key), that Issuer creates for a
// particular Credential Schema.
//
// Note! You should us WriteSchema, WriteDID, ...
// instead.
func BuildGetCredDefRequest(submitter, id string) ctx.Channel {
	return c2go.FindyBuildGetCredDefRequest(submitter, id)
}
