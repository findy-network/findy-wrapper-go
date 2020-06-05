// REMEMBER! const keywords are needed here

typedef void (*indy_handler)
	(indy_handle_t cmd_handle_, indy_error_t err);

typedef void (*indy_handler_str)
	(indy_handle_t cmd_handle_, indy_error_t err, const char *key);

typedef void (*indy_handler_str_str)
	(indy_handle_t cmd_handle_, indy_error_t err, const char *str1, const char *str2);

typedef void (*indy_handler_str_str_str)
	(indy_handle_t cmd_handle_, indy_error_t err, const char *str1, const char *str2, const char *str3);

typedef void (*indy_handler_str_str_ull)
	(indy_handle_t cmd_handle_, indy_error_t err, const char *str1, const char *str2, unsigned long long ull);

typedef void (*indy_handler_handle)
	(indy_handle_t cmd_handle_, indy_error_t err, indy_handle_t handle);

typedef void (*indy_handler_uint)
	(indy_handle_t cmd_handle_, indy_error_t err, unsigned int handle);

typedef void (*indy_handler_bool)
	(indy_handle_t cmd_handle_, indy_error_t err, bool handle);

typedef void (*indy_handler_handle_u32)
	(indy_handle_t cmd_handle_, indy_error_t err, indy_handle_t handle, indy_u32_t len);

typedef void (*indy_handler_u8ptr_u32)
	(indy_handle_t cmd_handle_, indy_error_t err, const indy_u8_t* encrypted_msg_raw, indy_u32_t encrypted_msg_len);

typedef void (*indy_handler_str_u8ptr_u32)
	(indy_handle_t cmd_handle_, indy_error_t err, const char *key, const indy_u8_t* encrypted_msg_raw, indy_u32_t encrypted_msg_len);

