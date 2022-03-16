# findy-wrapper-go

[![test](https://github.com/findy-network/findy-wrapper-go/actions/workflows/test.yml/badge.svg?branch=dev)](https://github.com/findy-network/findy-wrapper-go/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/findy-network/findy-wrapper-go/branch/dev/graph/badge.svg?token=2OPADTJQJ3)](https://codecov.io/gh/findy-network/findy-wrapper-go)

This is a Go wrapper for [indy-sdk](https://github.com/hyperledger/indy-sdk). It
wraps most of the **libindy**'s functions and data types, but it doesn't try to be
complete. We have written it incrementally and wrapped only those functions that
we need at the moment.

We haven't extended this wrapper with helper classes. It only tries to offer Go
programming interface above the existing libindy API. [findy-agent](https://github.com/findy-network/findy-agent)
extends this wrapper with SSI/DID abstractions and offers more high-level API
for SSI/DID development. For the most applications **findy-agent** is the one to
start playing with.

This wrapper offers one key feature addition to that it's written in Go: **it's
not dependent on running indy ledger.** The wrapper abstracts ledger with the
plug-in interface which allows to implement any storage technology to handle
needed persistence. Currently implemented ledger add-ons are:

- indy pool: data is saved into the indy ledger
- memory ledger: especially good for unit testing and caching
- file: data is saved into simple JSON file
- immudb: data is written to immutable database

## Get Started

Ubuntu 20.04 is preferred development environment but macOS is also an option.
Please make sure that Go and git are both installed and working properly.

### Linux and Ubuntu

This is the preferred way to build and use Findy Go wrapper.

1. Install libindy-dev: `make indy_to_debian`

   Check [indy-sdk installation instructions](https://github.com/hyperledger/indy-sdk/#installing-the-sdk) for more details.

2. Clone the repo: `git clone https://github.com/findy-network/findy-wrapper-go`
3. Run the tests to see everything is working properly : `make test`

### macOS

Because indy SDK won't offer proper distribution for OSX, we have written a
helper Bash script to perform installation. Follow these steps:

0. Install [Homebrew](https://brew.sh/) if it insn't already on your machine.
1. Clone the repo: `git clone https://github.com/findy-network/findy-wrapper-go`
2. Go to directory `./scripts/mac-libindy`:
   ```
   $ cd scripts/mac-libindy
   ```
3. Execute the installation script.
   ```
   $ ./install.sh
   ```
   **Or**, if you want to change the default installation location, enter it as
   a first argument for the script.
   ```
   $ ./install.sh /my/own/location
   ```
4. Follow the instructions of the `install.sh` i.e. **source the env.sh** which
   is generated to installation directory and is in your clipboard after successful
   installation.
   ```
   $ source /usr/local/opt/libindy/env.sh
   ```
5. Run the tests to see everything is OK: `make test`

The problem solving tip: `source env.sh` in your dev session.

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
2. Create a ledger pool with [indy CLI](https://github.com/bcgov/von-network/blob/master/docs/UsingVONNetwork.md#using-the-cli) on VON Network or if `findy-agent` is installed

   `findy-agent ledger pool create --name <pool_name> --genesis-txn-file genesis.txt`

   To test ledger connection you can give the following command:

   `findy-agent ledger pool ping --name <pool_name>`

3. Set environment variable: `export FINDY_POOL="FINDY_LEDGER,<pool_name>"`

   This will tell `anoncreds` tests to use `FINDY_LEDGER` plugin that's the
   add-on for real Indy ledger. You give the pool name for the ledger after
   the comma.

4. Run tests: `make test`.

   This will run all the pool tests towards the real ledger.

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
