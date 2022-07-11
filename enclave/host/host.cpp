#include <cstdio>
#include <cstring>
#include <fstream>
#include <iomanip>
#include <ios>
#include <iostream>
#include <memory>
#include <openenclave/bits/result.h>
#include <openenclave/host.h>
#include <string>
#include <sys/types.h>

#include "project_u.h"

oe_enclave_t *enclave = NULL;

int _create_enclave() {
  oe_result_t result;
  char enclave_path[] = "../build/enclave/enclave.signed";
  uint32_t flags = OE_ENCLAVE_FLAG_DEBUG;

  result = oe_create_project_enclave((char *)enclave_path, OE_ENCLAVE_TYPE_SGX,
                                     flags, NULL, 0, &enclave);
  if (result != OE_OK) {
    std::cerr << "oe_create_project_enclave(): result=" << result
              << " message=" << oe_result_str(result) << std::endl;
    if (enclave)
      oe_terminate_enclave(enclave);
    enclave = NULL;
  }
  return result;
}

// This function is called from inside the enclave(ocall)
bool host_store_data(const char *path, data_t *data) {
  printf("HOST: host_save_blob(..)\n");
  printf("HOST: PATH: %s, SIZE: %zu\n", path, data->size);
  printf("HOST BLOB: ");
  for (int i = 0; i < data->size; ++i) {
    printf("%d", data->blob[i]);
  }
  printf("\n");
  std::fstream f{path, std::ios_base::binary | std::ios_base::out};

  if (!f.is_open()) {
    printf("HOST: couldn't open the file\n");
    return false;
  }

  f.write((char *)&data->blob[0], data->size);
  bool ret = f.good();
  f.close();

  return ret;
}

bool host_load_data(const char *path, data_t *sealed_data) {
  std::ifstream f{path, std::ios_base::binary | std::ios_base::in |
                            std::ios_base::ate};
  bool ret = false;
  data_t s_data;

  if (f.is_open()) {
    s_data.size = f.tellg();
    printf("HOST: blob size: %zu\n", s_data.size);
    if (s_data.size > 0) {
      s_data.blob = new uint8_t[s_data.size];
      f.seekg(0);
      f.read((char *)s_data.blob, s_data.size);
      ret = f.good();
    }
  } else {
    printf("HOST: File not found: %s\n", path);
  }
  f.close();
  sealed_data->blob = s_data.blob;
  sealed_data->size = s_data.size;

  return ret;
}

// During the /register/finish the server calls this function to store and
// seal the public key
int host_store_ecc_pk(unsigned char *x, unsigned char *y) {
  bool ret = false;
  if (_create_enclave() == OE_OK) {
    enclave_store_ecc_pk(enclave, &ret, x, y);
  }
  return ret;
}

unsigned char *host_create_nonce(unsigned int len) {
  unsigned char *nonce = new unsigned char[len];
  std::memset(nonce, 0, len);
  printf("host_create_nonce()\n");

  if (_create_enclave() == OE_OK) {
    printf("create enclave\n");
    enclave_create_nonce(enclave, nonce, len, NULL);
    printf("nonce: %s\n", nonce);
  }
  return nonce;
}

unsigned char *host_create_nonce_hash(unsigned char *hash, unsigned int len) {
  data_t hash_data{hash, len};
  unsigned char *nonce = new unsigned char[len + 32];
  std::memset(nonce, 0, len + 32);
  printf("host_create_nonce_hash()\n");

  if (_create_enclave() == OE_OK) {
    printf("create enclave\n");
    enclave_create_nonce(enclave, nonce, len + 32, &hash_data);
    printf("nonce: %s\n", nonce);
  }
  return nonce;
}

int host_verify_secp256r1_sig(unsigned char *msg, unsigned int msg_len,
                              unsigned char *sig, unsigned int sig_len) {
  printf("host_verify_secp256r1_sig()\n");

  bool enc_ret = false;
  data_t data{msg, msg_len};
  data_t pk_sealed{NULL, 0};

  if ((enc_ret = host_load_data("att_pk.bin", &pk_sealed))) {
    printf("HOST: loaded the sealed public key\n");
    printf("HOST: sealed public key size: %zu\n", pk_sealed.size);
    if (_create_enclave() == OE_OK) {
      printf("HOST: Enclave created\n");
      enclave_verify_secp256r1_sig(enclave, &enc_ret, &data, &pk_sealed, sig,
                                   sig_len);
    }
  } else {
    printf("HOST: ERROR: Couldn't load the sealed public key\n");
  }
  delete[] pk_sealed.blob;
  printf("HOST: return value is: %d\n", enc_ret);
  return enc_ret;
}

int host_gen_secp256k1_keys() {
  int ret = 0;
  printf("host_create_wallet_keys()\n");

  if (_create_enclave() == OE_OK) {
    printf("HOST: Created enclave\n");
    enclave_gen_secp256k1_keys(enclave, &ret);
  }
  return ret;
}

point *host_get_pubkey() {
  point *pubkey = new point{0, 0};
  printf("host_get_pubkey()\n");
  data_t sealed_data{NULL, 0};
  if (!host_load_data("wallet_pubkey.bin", &sealed_data)) {
    if (host_gen_secp256k1_keys() != OE_OK) {
      printf("HOST: ERROR: An error occured during the generation of the "
             "keypair\n");
      goto exit;
    }

    if (_create_enclave() == OE_OK) {
      printf("HOST: ENCLAVE_CREATED\n");
      enclave_get_pubkey(enclave, pubkey, &sealed_data);
    }
  }
exit:
  return pubkey;
}

unsigned char *host_sign_secp256k1(unsigned char *hash, unsigned int hash_len) {
  printf("host_sign_secp256k1()\n");
  data_t sealed_bin;
  data_t sig_data;
  if (!host_load_data("wallet_privkey.bin", &sealed_bin)) {
    if (host_gen_secp256k1_keys() != OE_OK) {
      printf("HOST: ERROR: An error occured during the generation of the "
             "keypair\n");
      goto exit;
    }

    data_t hash_data{hash, hash_len};
    if (_create_enclave() == OE_OK) {
      printf("HOST: ENCLAVE_CREATED\n");
      enclave_sign_sha256(enclave, &hash_data, &sealed_bin, &sig_data);
    }
  }
exit:
  return sig_data.blob;
}

// Only for testing
void test_sign_secp256k1() {
  printf("host_test_secp256k1()\n");
  data_t sealed_privkey;
  data_t sealed_pubkey;
  data_t signed_data{NULL, 0};
  data_t hash{(unsigned char *)std::malloc(32), 32};
  char *str = (char *)"Test";

  // Load the sealed private key
  if (!host_load_data("wallet_privkey.bin", &sealed_privkey)) {
    printf("HOST: Generating keypair\n");
    if (host_gen_secp256k1_keys() != OE_OK) {
      printf("HOST: ERROR: An error occured during the generation of the "
             "keypair\n");
      return;
    }
  }

  // Load the sealed public key
  host_load_data("wallet_pubkey.bin", &sealed_pubkey);

  // Generate random value
  if (_create_enclave() == OE_OK) {
    printf("HOST: ENCLAVE_CREATED\n");
    enclave_hash(enclave, str, hash.blob);
    printf("HOST: Hashed value\n");
    enclave_sign_sha256(enclave, &hash, &sealed_privkey, &signed_data);
    printf("HOST: Signed data\n");
    bool ret = false;
    enclave_verify_secp256k1_sig(enclave, &ret, &sealed_pubkey, &signed_data,
                                 &hash);
  }
}
