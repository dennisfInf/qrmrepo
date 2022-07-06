unsigned char* host_create_nonce(unsigned int len);
// return 1 on success
int host_verify_secp256r1_sig(unsigned char *msg, unsigned int msg_len,
                              unsigned char *sig, unsigned int sig_len);
// return 1 on success
int host_store_ecc_pk(unsigned char *x, unsigned char *y);
