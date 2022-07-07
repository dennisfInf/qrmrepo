unsigned char *host_create_nonce(unsigned int len);
unsigned char *host_create_nonce_hash(unsigned char *hash, unsigned int len);
// return 1 on success
int host_verify_secp256r1_sig(unsigned char *msg, unsigned int msg_len,
                              unsigned char *sig, unsigned int sig_len);
// return 1 on success
int host_store_ecc_pk(unsigned char *x, unsigned char *y);
int host_gen_secp256k1_keys();

struct point {
  unsigned long x, y;
};

struct point* host_get_pubkey();
unsigned char *host_sign_secp256k1(unsigned char *hash, unsigned int hash_len);
