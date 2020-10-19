# findy-wrapper-go

This is a Go wrapper for [indy-sdk](https://github.com/hyperledger/indy-sdk). It
wraps most of the **libindy**'s functions and data types, but it doesn't try to be
complete. We have written it incrementally and wrapped only those functions that
we need at the moment.

We haven't extended this wrapper with helper classes. It only tries to offer Go
programming interface above the existing libindy API. [findy-agent](https://github.com/findy-network/findy-agent) 
extends this API with SSI/DID abstractions and offers more high-level API for
SSI/DID development. For the most applications **findy-agent** is the suitable
framework to start. 

## Get Started

1. [Install](https://github.com/hyperledger/indy-sdk/#installing-the-sdk) libindy-dev.
2. Clone the repo: `git clone https://github.com/findy-network/findy-go`
3. Install needed Go packages: `make deps`
4. Build the package: `make build`

If build system cannot find indy libs and headers, set following environment 
variables:

```
export CGO_CFLAGS="-I/<path_to_>/indy-sdk/libindy/include"
export CGO_LDFLAGS="-L/<path_to_>/indy-sdk/libindy/target/debug"
```

## Development Without The Ledger

The findy-go package includes memory implementation of the ledger transaction
interface. If pool names are `FINDY_MEM_LEDGER` or `FINDY_ECHO_LEGDER` the
package uses its internal implementation of the ledger pool and lets you
perform all of its functions.

```go
	r := <-pool.OpenLedger("FINDY_MEM_LEDGER")
	if r.Err() != nil {
		return r.Err()
	}
	pool := r.Handle()
```

## Run Tests With Indy Ledger

1. [Install and start ledger](https://github.com/bcgov/von-network/blob/master/docs/UsingVONNetwork.md#building-and-starting)
2. Create a ledger pool with [indy CLI](https://github.com/bcgov/von-network/blob/master/docs/UsingVONNetwork.md#using-the-cli)  on VON Network or if `findy-agent` is installed

   ```findy-agent create cnx -pool <pool_name> -txn genesis.txt```
3. Set environment variable: `export FINDY_POOL=<pool_name>`
4. Run tests: `make test`

## Documentation

The wrapper includes minimal Go documentation. If you need more information
about indy SDK, we suggest that you would read the original 
[indy SDK documentation](https://hyperledger-indy.readthedocs.io/projects/sdk/en/latest/docs/index.html).

## Naming Conventions

It follows the same sub-package structure as libindy.
Also, function names are the same, but it respects Go idioms.

Original indy SDK function calls in C:
```C
indy_key_for_did(...)
...
indy_key_for_local_did(..)
```
Same functions but with Go wrapper:
```go
r := <-did.Key(pool, wallet, didStr)
...
r = <-did.LocalKey(wallet, didStr)
```

As you can see, the subjects that would exist in the function names aren't
repeated if they are in the package names already.

## Return Types

The return type is Channel, which transfers `dto.Result` Go structures. These
structures work similarly to C unions, which means that we can use one type for
all of the wrapper functions. To access the actual data, you have to know the
actual type. `pool.OpenLedger()` returns `Handle` type.

```go
	r = <-pool.OpenLedger("FINDY_MEM_LEDGER")
	assert.NoError(t, r.Err())
	h2 := r.Handle()
	assert.Equal(t, h2, -1)
```

`did.CreateAndStore()` returns two strings: `did` and `verkey`. Please note the
use of `did.Did` struct instead of the JSON string.

```go
	r := <-did.CreateAndStore(w, did.Did{Seed: ""})
	assert.NoError(t, r.Err())
	did := r.Str1()
	vk := r.Str2()
```

When a null string is needed for an argument, the predefined type must be used.

```go
	r = <-annoncreds.IssuerCreateAndStoreCredentialDef(w1, w1DID, scJSON,
		"MY_FIRMS_CRED_DEF", findy.NullString, findy.NullString)
```

## Error Values

The Go error value can be retrieved with `dto.Result.Err()` which returns Go
`error`.