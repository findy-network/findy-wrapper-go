//
// Autogenerated from indy-sdk main header files (indy_*.h)
//

#include <stdio.h>
#include "findy_glue.h"
#include "findy_callback_types.h"

// cgo generated
#include "_cgo_export.h"


indy_error_t findy_is_pairwise_exists(indy_handle_t command_handle, indy_handle_t wallet_handle, char *their_did ) {
	indy_error_t err = indy_is_pairwise_exists(command_handle, wallet_handle, their_did, (indy_handler_bool)boolHandler);
	if (err != Success) {
		handler(command_handle, err);
	}
	return err;
}

indy_error_t findy_create_pairwise(indy_handle_t command_handle, indy_handle_t wallet_handle, char *their_did, char *my_did, char *metadata ) {
	indy_error_t err = indy_create_pairwise(command_handle, wallet_handle, their_did, my_did, metadata, (indy_handler)handler);
	if (err != Success) {
		handler(command_handle, err);
	}
	return err;
}

indy_error_t findy_list_pairwise(indy_handle_t command_handle, indy_handle_t wallet_handle ) {
	indy_error_t err = indy_list_pairwise(command_handle, wallet_handle, (indy_handler_str)strHandler );
	if (err != Success) {
		handler(command_handle, err);
	}
	return err;
}

indy_error_t findy_get_pairwise(indy_handle_t command_handle, indy_handle_t wallet_handle, char *their_did ) {
	indy_error_t err = indy_get_pairwise(command_handle, wallet_handle, their_did, (indy_handler_str)strHandler );
	if (err != Success) {
		handler(command_handle, err);
	}
	return err;
}

indy_error_t findy_set_pairwise_metadata(indy_handle_t command_handle, indy_handle_t wallet_handle, char *their_did, char *metadata ) {
	indy_error_t err = indy_set_pairwise_metadata(command_handle, wallet_handle, their_did, metadata, (indy_handler)handler);
	if (err != Success) {
		handler(command_handle, err);
	}
	return err;
}

