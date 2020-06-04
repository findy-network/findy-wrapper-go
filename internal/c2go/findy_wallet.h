extern indy_error_t findy_create_wallet(indy_handle_t command_handle, char*config, char*credentials);
extern indy_error_t findy_open_wallet(indy_handle_t command_handle, char*config, char*credentials);
extern indy_error_t findy_export_wallet(indy_handle_t command_handle, indy_handle_t wallet_handle, char*export_config_json);
extern indy_error_t findy_import_wallet(indy_handle_t command_handle, char*config, char*credentials, char*import_config_json);
extern indy_error_t findy_close_wallet(indy_handle_t command_handle, indy_handle_t wallet_handle);
extern indy_error_t findy_delete_wallet(indy_handle_t command_handle, char*config, char*credentials);
extern indy_error_t findy_generate_wallet_key(indy_handle_t command_handle, char *config);
