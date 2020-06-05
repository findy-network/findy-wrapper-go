extern indy_error_t findy_is_pairwise_exists(indy_handle_t command_handle, indy_handle_t wallet_handle, char *their_did);
extern indy_error_t findy_create_pairwise(indy_handle_t command_handle, indy_handle_t wallet_handle, char *their_did, char *my_did, char *metadata);
extern indy_error_t findy_list_pairwise(indy_handle_t command_handle, indy_handle_t wallet_handle);
extern indy_error_t findy_get_pairwise(indy_handle_t command_handle, indy_handle_t wallet_handle, char *their_did);
extern indy_error_t findy_set_pairwise_metadata(indy_handle_t command_handle, indy_handle_t wallet_handle, char *their_did, char *metadata);
