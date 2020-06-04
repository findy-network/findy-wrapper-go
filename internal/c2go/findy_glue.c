#include <stdio.h>

#include "findy_glue.h"

// cgo generated
#include "_cgo_export.h"

#include "findy_callback_types.h"

char* findy_null_string = NULL;

char* findy_get_current_error() {
	static char* errorJson;
	indy_get_current_error((const char **)&errorJson);
	return errorJson;
}

indy_error_t findy_create_pool_ledger_config(indy_handle_t cmd_handle, char *config_name, char *config) {
	indy_error_t err = indy_create_pool_ledger_config(cmd_handle, config_name, config, (indy_handler)handler );
	if (err != Success) {
		handler(cmd_handle, err);
	}
	return err;
}

indy_error_t findy_open_pool_ledger(indy_handle_t cmd_handle, char *config_name, char *config) {
	indy_error_t err = indy_open_pool_ledger(cmd_handle, config_name, config, (indy_handler_handle)handleHandler );
	if (err != Success) {
		handleHandler(cmd_handle, err, 0);
	}
	return err;
}

indy_error_t findy_close_pool_ledger(indy_handle_t cmd_handle, indy_handle_t handle) {
	indy_error_t err = indy_close_pool_ledger(cmd_handle, handle, (indy_handler)handler );
	if (err != Success) {
		handler(cmd_handle, err);
	}
	return err;
}

indy_error_t findy_list_pools(indy_handle_t cmd_handle) {
	indy_error_t err = indy_list_pools(cmd_handle, (indy_handler_str)strHandler );
	if (err != Success) {
		strHandler(cmd_handle, err, "ERROR");
	}
	return err;
}

indy_error_t findy_set_protocol_version(indy_handle_t cmd_handle, indy_u64_t protocol_version) {
	indy_error_t err = indy_set_protocol_version(cmd_handle, protocol_version, (indy_handler)handler );
	if (err != Success) {
		handler(cmd_handle, err);
	}
	return err;
}
