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
  printf("HOST: PATH: %s, SIZE: %d\n", path, data->size);
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
  std::fstream f{path, std::ios_base::binary | std::ios_base::in |
                           std::ios_base::ate};
  bool ret = false;
  data_t s_data;

  if (f.is_open()) {
    s_data.size = f.tellg();
    printf("HOST: blob size: %d\n", s_data.size);
    if (s_data.size > 0) {
      s_data.blob = new uint8_t[s_data.size];
      f.seekg(0);
      f.read((char *)s_data.blob, s_data.size);
      ret = f.good();
    }
  }
  f.close();

  printf("HOST: blob: ");
  for (int i = 0; i < s_data.size; ++i) {
    printf("%d", s_data.blob[i]);
  }
  printf("\n");

  printf("HOST: FINISH\n");
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
    enclave_create_nonce(enclave, nonce, len);
    printf("nonce: %s\n", nonce);
  }
  printf("Finish");
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
    printf("HOST: sealed public key size: %d\n", pk_sealed.size);
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

point* host_get_pubkey() {
	printf("host_get_pubkey()\n");
	data_t sealed_data{NULL, 0};
	host_load_data("wallet_pubkey.bin", &sealed_data);
	point *pubkey = new point{0,0};
	printf("HOST: CREATE ENCLAVE, %s\n", sealed_data.blob);
	if(_create_enclave() == OE_OK) {
		printf("HOST: ENCLAVE_CREATED\n");
		enclave_get_pubkey(enclave, pubkey, &sealed_data);
	}
	return pubkey;
}
