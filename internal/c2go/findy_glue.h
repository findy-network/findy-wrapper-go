
#include "indy/indy_core.h"

extern char* findy_null_string;

extern char* findy_get_current_error();

// regex helpers
//		,\s*$			remove not needed line breaks, not working in IDE
//		^\s+			remove long space line starts
//						remove comments
//		^\s*$			remove empty lines
//		^\s*///.*$		remove comment lines

// MARK: wallet API
#include "findy_wallet.h"

// MARK: pool API
// pool api, note reverse-engineered from cli
extern indy_error_t findy_create_pool_ledger_config(indy_handle_t command_handle, char *config_name, char *config);
// connect in cli, here open, in CLI has extra parameters like timeout and protocol which has own api in indy-sdk
extern indy_error_t findy_open_pool_ledger(indy_handle_t command_handle, char *config_name, char *config);
// disconnect in cli, here close
extern indy_error_t findy_close_pool_ledger(indy_handle_t command_handle, indy_handle_t handle);
extern indy_error_t findy_list_pools(indy_handle_t command_handle);
extern indy_error_t findy_set_protocol_version(indy_handle_t cmd_handle, indy_u64_t protocol_version);

#include "findy_ledger.h"
#include "findy_crypto.h"
#include "findy_anoncreds.h"
#include "findy_blob_storage.h"
#include "findy_pairwise.h"
#include "findy_did.h"