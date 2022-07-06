#include "project_t.h"
#include <algorithm>
#include <cstddef>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <mbedtls/asn1.h>
#include <mbedtls/bignum.h>
#include <mbedtls/ctr_drbg.h>
#include <mbedtls/ecdsa.h>
#include <mbedtls/ecp.h>
#include <mbedtls/entropy.h>
#include <mbedtls/sha256.h>
#include <memory>
#include <openenclave/bits/result.h>
#include <openenclave/bits/types.h>
#include <openenclave/enclave.h>
#include <openenclave/log.h>
#include <openenclave/seal.h>
#include <stdio.h>
#include <sys/types.h>

unsigned char *tmp_nonce = NULL;

#define PATH_ATTESTATION_PUBLIC_KEY "att_pk.bin"
#define PATH_NONCE "nonce.bin"

class Seal {
public:
  Seal(oe_seal_policy_t &&policy, const char *path);
  oe_result_t seal_store(data_t *data) const;
  bool store(data_t *data) const;
  oe_result_t seal(data_t *data, uint8_t **blob, size_t *blob_size) const;
  oe_result_t unseal(data_t *sealed, uint8_t **data, size_t *data_size) const;
  ~Seal();

protected:
  const char *path;
  const static size_t SEAL_SETTING_COUNT = 1;
  oe_seal_policy_t seal_policy;
  oe_seal_setting_t *seal_settings;

  // opt_size can't be 0 in this version of OpenEnclave
  uint8_t *opt_msg = (uint8_t *)"1";
  size_t opt_size = 1;
};

Seal::Seal(oe_seal_policy_t &&policy, const char *path)
    : seal_policy{policy}, path{path} {
  this->seal_settings = new oe_seal_setting_t[Seal::SEAL_SETTING_COUNT]{
      OE_SEAL_SET_POLICY(this->seal_policy)};
}

Seal::~Seal() { delete this->seal_settings; }

oe_result_t Seal::seal_store(data_t *data) const {
  oe_result_t ret = OE_FAILURE;
  data_t sealed{NULL, 0};

  if ((ret = this->seal(data, &sealed.blob, &sealed.size)) == OE_OK) {
    if (!this->store(&sealed)) {
      ret = OE_FAILURE;
    }
  }
  return OE_OK;
}

bool Seal::store(data_t *data) const {
  bool ret = false;
  host_store_data(&ret, this->path, data);
  return ret;
}

oe_result_t Seal::seal(data_t *data, uint8_t **blob, size_t *blob_size) const {
  printf("ENCLAVE: Storage::seal(..)\n");

  oe_result_t ret =
      oe_seal(NULL, this->seal_settings, this->SEAL_SETTING_COUNT, data->blob,
              data->size, this->opt_msg, this->opt_size, blob, blob_size);

  if (*blob_size > UINT32_MAX) {
    printf("blob_size is too large to fit into an unsigned int\n");
    ret = OE_OUT_OF_MEMORY;
  }
  return ret;
}

oe_result_t Seal::unseal(data_t *sealed, uint8_t **data,
                         size_t *data_size) const {
  printf("ENCLAVE: Seal::unseal()\n");
  printf("\tsealed length: %d\n", sealed->size);
  return oe_unseal(sealed->blob, sealed->size, this->opt_msg, this->opt_size,
                   data, data_size);
}

void enclave_test() {
  printf("ENCLAVE: enclave_test()\n");
  uint8_t *data = (uint8_t *)"dd";
  size_t data_size = 2;
  uint8_t *blob;
  size_t blob_size;
  // oe_result_t ret = seal_data(data, 2, blob, &blob_size);
  // SealStorage storage{OE_SEAL_POLICY_PRODUCT};
  // auto ret = storage.seal(data, &data_size);
  auto ret = OE_OK;

  Seal seal{OE_SEAL_POLICY_PRODUCT, "enclave.bin"};

  switch (ret) {
  case OE_OK:
    printf("ENCLAVE: Data successfully sealed\n");
    break;
  case OE_INVALID_PARAMETER:
    printf("ENCLAVE: Invalid parameter\n");
    break;
  case OE_UNSUPPORTED:
    printf("ENCLAVE: unsupported\n");
    break;
  case OE_OUT_OF_MEMORY:
    printf("ENCLAVE: Out of memory\n");
    break;
  case OE_CRYPTO_ERROR:
    printf("ENCLAVE: Crypto Error\n");
    break;
  default:
    printf("ENCLAVE: Unkown error\n");
  }
}

void enclave_hash(char *str, unsigned char hash[]) {
  mbedtls_sha256_context ctx;
  mbedtls_sha256_init(&ctx);
  mbedtls_sha256_starts_ret(&ctx, false);

  mbedtls_sha256_update_ret(&ctx, (unsigned char *)str, strlen(str));
  mbedtls_sha256_finish_ret(&ctx, hash);
  mbedtls_sha256_free(&ctx);
}

oe_result_t unseal_data(const uint8_t *blob, size_t blob_size, uint8_t **data,
                        size_t *data_size) {
  printf("ENCLAVE: unseal_data()\n");

  const oe_seal_setting_t seal_settings[] = {
      OE_SEAL_SET_POLICY(OE_SEAL_POLICY_PRODUCT),
  };
  return oe_unseal(blob, blob_size, NULL, 0, data, data_size);
}

// x and y must be binaries with a length of exactly 32 bytes
bool enclave_store_ecc_pk(unsigned char *x, unsigned char *y) {
  bool ret = false;
  mbedtls_ecdsa_context ctx_verify;
  mbedtls_ecp_group grp;
  mbedtls_ecdsa_init(&ctx_verify);
  mbedtls_ecp_group_init(&grp);

  if ((ret = mbedtls_ecp_group_load(&ctx_verify.grp,
                                    MBEDTLS_ECP_DP_SECP256R1)) == 0) {
    printf("ENCLAVE: Loaded the ecc-group\n");

    // Create the binary which mbedtls_ecp_point_read_binary(...) needs
    // 1. This binary has to be exactly 65 bytes
    // 2. In this case the points aren't compressed, so that the first entry
    // in the binary has to be 0x04
    unsigned char *tmp = (unsigned char *)std::malloc(65);
    tmp[0] = 0x04;
    std::memcpy(tmp + 1, x, 32);
    std::memcpy(tmp + 33, y, 32);

    data_t data{tmp, 65};
    Seal seal{OE_SEAL_POLICY_PRODUCT, PATH_ATTESTATION_PUBLIC_KEY};
    ret = seal.seal_store(&data);
  } else {
    printf("ENCLAVE: failed: !mbedtls_ecp_group_load returned %d\n", ret);
  }

  // Free memory allocations
  mbedtls_ecp_group_free(&grp);
  mbedtls_ecdsa_free(&ctx_verify);
  return ret;
}

// Generate a nonce and write it to out
// Store the nonce
void enclave_create_nonce(unsigned char *out, uint8_t len) {
  printf("ENCLAVE: enclave_create_nonce()\n");

  int ret = 1;
  int exit_code = EXIT_FAILURE;
  mbedtls_ctr_drbg_context ctr_drbg;
  mbedtls_entropy_context entropy;

  mbedtls_ctr_drbg_init(&ctr_drbg);
  mbedtls_entropy_init(&entropy);

  ret =
      mbedtls_ctr_drbg_seed(&ctr_drbg, mbedtls_entropy_func, &entropy, NULL, 0);
  if (ret == 0) {
    mbedtls_ctr_drbg_set_prediction_resistance(&ctr_drbg,
                                               MBEDTLS_CTR_DRBG_PR_OFF);
    ret = mbedtls_ctr_drbg_random(&ctr_drbg, out, 32);
    if (ret == 0) {
      exit_code = EXIT_SUCCESS;
    } else {
      printf("failed in mbedtls_ctr_drbg_random: %d\n", ret);
    }
  } else {
    printf("failed in mbedtls_ctr_drbg_seed: %d\n", ret);
  }

  mbedtls_ctr_drbg_free(&ctr_drbg);
  mbedtls_entropy_free(&entropy);

  Seal seal{OE_SEAL_POLICY_PRODUCT, PATH_NONCE};
  data_t data{out, len};
  seal.seal_store(&data);
}

// This function is used in webauthn to check the signature of the passed data
bool enclave_verify_secp256r1_sig(data_t *data, data_t *pk, unsigned char *sig,
                                  unsigned int sig_len) {
  printf("called enclave_verify_secp256r1_sig()\n");

  // Print all the inputs for debug reasons
  printf("ENCLAVE: data: ");
  for (int i = 0; i < data->size; ++i) {
    printf("%02x", data->blob[i]);
  }
  printf("\nENCLAVE: sig: ");
  for (int i = 0; i < sig_len; ++i) {
    printf("%02x", sig[i]);
  }
  printf("\n");

  int ret = 0;
  unsigned char *hash = (unsigned char *)std::malloc(32);
  mbedtls_ecdsa_context ctx_verify;
  mbedtls_ecp_group grp;

  mbedtls_ecdsa_init(&ctx_verify);
  mbedtls_ecp_group_init(&grp);

  Seal seal{OE_SEAL_POLICY_PRODUCT, PATH_ATTESTATION_PUBLIC_KEY};
  uint8_t *pk_bin = NULL;
  size_t pk_bin_size = 0;

  if (seal.unseal(pk, &pk_bin, &pk_bin_size) == OE_OK) {
    printf("ENCLAVE: Unsealed the public key\n");

    // Compute sha256 sum
    printf("ENCLAVE: Hashing the data\n");
    mbedtls_sha256_ret(data->blob, data->size, hash, 0);

    if ((ret = mbedtls_ecp_group_load(&ctx_verify.grp,
                                      MBEDTLS_ECP_DP_SECP256R1)) == 0) {
      printf("ENCLAVE: Loaded the ecc-group\n");

      printf("ENCLAVE: Read x and y for ecc\n");

      // Read the points of the public key from the created binary
      // pk->size must be exactly 65 bytes
      if (mbedtls_ecp_point_read_binary(&ctx_verify.grp, &ctx_verify.Q, pk_bin,
                                        pk_bin_size) == 0) {
        printf("ENCLAVE: Successfully read the point from the binary:\n");

        // Verify the correctness of the public key
        if (mbedtls_ecp_check_pubkey(&ctx_verify.grp, &ctx_verify.Q) == 0) {
          printf("ENCLAVE: The public key is valid\n");
          // Verify the signature
          if ((ret = mbedtls_ecdsa_read_signature(&ctx_verify, hash, 32, sig,
                                                  sig_len)) == 0) {
            printf("ENCLAVE: The signature is valid\n");
          } else {
            printf("ENCLAVE: ERROR: The signature is invalid\n");
            printf("ENCLAVE: ERROR: signature: %s", sig);
          }
        } else {
          printf("ENCLAVE: ERROR: The public key is invalid\n");
        }
      } else {
        printf("ENCLAVE: ERROR: Couldn't read the ec public key point from "
               "binary\n");
      }
    } else {
      printf("ENCLAVE: ERROR: !mbedtls_ecp_group_load returned %d\n", ret);
    }
  } else {
    printf("ENCLAVE: ERROR: Couldn't unseal the public key binary\n");
  }

  // Free memory allocations
  mbedtls_ecp_group_free(&grp);
  mbedtls_ecdsa_free(&ctx_verify);

  std::free(hash);

  return ret == 0;
}
