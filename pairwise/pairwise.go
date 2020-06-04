/*
Package pairwise is corresponding Go package for libindy's pairwise namespace.
We suggest that you read indy SDK documentation for more information.
*/
package pairwise

import (
	"github.com/optechlab/findy-go/internal/c2go"
	"github.com/optechlab/findy-go/internal/ctx"
)

// Exists returns true if pairwise exists in the wallet.
func Exists(wallet int, theirDID string) ctx.Channel {
	return c2go.FindyIsPairwiseExists(wallet, theirDID)
}

// Create creates a new pairwise to wallet with the argumets.
func Create(wallet int, theirDID, myDID, metaData string) ctx.Channel {
	return c2go.FindyCreatePairwise(wallet, theirDID, myDID, metaData)
}

// List lists lists all of the pairwise in the wallet.
func List(wallet int) ctx.Channel {
	return c2go.FindyListPairwise(wallet)
}

// Get returns a pairwise data with their DID.
func Get(wallet int, theirDID string) ctx.Channel {
	return c2go.FindyGetPairwise(wallet, theirDID)
}

// SetMeta sets the meta data for their DID.
func SetMeta(wallet int, theirDID, meta string) ctx.Channel {
	return c2go.FindySetPairwiseMetadata(wallet, theirDID, meta)
}
