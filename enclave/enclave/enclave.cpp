#include "project_args.h"
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
#define PATH_WALLET_PUBKEY "wallet_pubkey.bin"
#define PATH_WALLET_PRIVKEY "wallet_privkey.bin"

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
  printf("\tsealed length: %zu\n", sealed->size);
  return oe_unseal(sealed->blob, sealed->size, this->opt_msg, this->opt_size,
                   data, data_size);
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
void enclave_create_nonce(unsigned char *out, uint8_t len, data_t *opt) {
  printf("ENCLAVE: enclave_create_nonce()\n");

  int ret = 1;
  int exit_code = EXIT_FAILURE;
  size_t nonce_size = 32;
  mbedtls_ctr_drbg_context ctr_drbg;
  mbedtls_entropy_context entropy;

  mbedtls_ctr_drbg_init(&ctr_drbg);
  mbedtls_entropy_init(&entropy);

  ret =
      mbedtls_ctr_drbg_seed(&ctr_drbg, mbedtls_entropy_func, &entropy, NULL, 0);
  if (ret == 0) {
    mbedtls_ctr_drbg_set_prediction_resistance(&ctr_drbg,
                                               MBEDTLS_CTR_DRBG_PR_OFF);
    ret = mbedtls_ctr_drbg_random(&ctr_drbg, out, nonce_size);
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
  if (opt != NULL && opt->size > 0) {
    size_t challenge_size = nonce_size + opt->size;
    unsigned char *challenge = new unsigned char[challenge_size];

    std::memcpy(challenge, out, nonce_size);
    std::memcpy(challenge + nonce_size, opt->blob, opt->size);

    data_t data{challenge, challenge_size};
    seal.seal_store(&data);
    return;
  }
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

void enclave_gen_secp256k1_keys(int *ret) {
  mbedtls_ecdsa_context ctx;
  mbedtls_entropy_context entropy;
  mbedtls_ecdsa_init(&ctx);
  mbedtls_entropy_init(&entropy);

  // Generate the keypair
  if ((*ret = mbedtls_ecdsa_genkey(&ctx, MBEDTLS_ECP_DP_SECP256K1,
                                   mbedtls_entropy_func, &entropy)) != 0) {
    printf("ENCLAVE: ERROR: mbedtls_ecdsa_genkey() returned with %d\n", *ret);
  } else {
    printf("ENCLAVE: x: %lu, y: %lu\n", *ctx.Q.X.p, *ctx.Q.Y.p);
    // Write the public key to a binary
    size_t pkey_binary_len = 65;
    unsigned char *pkey_binary = new unsigned char[pkey_binary_len];
    if ((*ret = mbedtls_ecp_point_write_binary(
             &ctx.grp, &ctx.Q, MBEDTLS_ECP_PF_UNCOMPRESSED, &pkey_binary_len,
             pkey_binary, pkey_binary_len)) != 0) {
      printf(
          "ENCLAVE: ERROR: mbedtls_ecp_point_write_binary() returned with %d\n",
          *ret);
    } else {
      // Seal and store the public key
      data_t pubkey_data{pkey_binary, pkey_binary_len};
      Seal seal{OE_SEAL_POLICY_PRODUCT, PATH_WALLET_PUBKEY};
      if ((*ret = seal.seal_store(&pubkey_data)) != 0) {
        printf("ENCLAVE: ERROR: seal_store() returned with %d\n", *ret);
      } else {
        // Write the private key to a binary
        size_t privkey_binary_len = mbedtls_mpi_size(&ctx.d);
        unsigned char *privkey_binary = new unsigned char[privkey_binary_len];

        if ((*ret = mbedtls_mpi_write_binary(&ctx.d, privkey_binary,
                                             privkey_binary_len)) != 0) {
          printf(
              "ENCLAVE: ERROR: mbedtls_mpi_write_binary() returned with %d\n",
              *ret);
        } else {
          // Seal and store the private key
          data_t privkey_data{privkey_binary, privkey_binary_len};
          Seal seal{OE_SEAL_POLICY_PRODUCT, PATH_WALLET_PRIVKEY};
          if ((*ret = seal.seal_store(&privkey_data)) != 0) {
            printf("ENCLAVE: ERROR: seal_store() returned with %d\n", *ret);
          }
        }
        delete[] privkey_binary;
      }
    }
    delete[] pkey_binary;
  }

  // Seal and store the private key

  mbedtls_ecdsa_free(&ctx);
  mbedtls_entropy_free(&entropy);
}

void enclave_get_pubkey(point *pubkey, data_t *sealed_data) {
  printf("ENCLAVE: enclave_get_pubkey()\n");
  int ret = 0;
  mbedtls_ecp_group grp;
  mbedtls_ecp_point p;
  mbedtls_ecp_group_init(&grp);
  mbedtls_ecp_point_init(&p);
  mbedtls_ecp_group_load(&grp, MBEDTLS_ECP_DP_SECP256K1);
  unsigned char *data = NULL;
  size_t data_size = 0;
  Seal seal{OE_SEAL_POLICY_PRODUCT, PATH_WALLET_PUBKEY};
  if (seal.unseal(sealed_data, &data, &data_size) == 0) {
    if ((ret = mbedtls_ecp_point_read_binary(&grp, &p, data, data_size)) != 0) {
      printf(
          "ENCLAVE: ERROR: mbedtls_ecp_point_read_binary() returned with %d\n",
          ret);
    } else {
      pubkey->x = *p.X.p;
      pubkey->y = *p.Y.p;
      printf("ENCLAVE: x: %lu, y: %lu\n", *p.X.p, *p.Y.p);
    }
  }
}

void enclave_sign_sha256(data_t *hash_data, data_t *sealed_bin,
                         data_t *sig_data) {
  if (hash_data == NULL || sealed_bin == NULL) {
    return;
  }

  unsigned char *data = NULL;
  size_t data_size = 0;
  Seal seal_bin{OE_SEAL_POLICY_PRODUCT, PATH_WALLET_PRIVKEY};
  seal_bin.unseal(sealed_bin, &data, &data_size);
  mbedtls_mpi priv;
  mbedtls_entropy_context entropy;
  mbedtls_mpi_init(&priv);

  // Read the private key
  mbedtls_mpi_read_binary(&priv, data, data_size);
  mbedtls_ecdsa_context ctx;
  mbedtls_ecp_group grp;
  mbedtls_ecp_point p;
  mbedtls_entropy_init(&entropy);
  mbedtls_ecdsa_init(&ctx);
  mbedtls_ecp_group_init(&grp);
  mbedtls_ecp_point_init(&p);
  mbedtls_ecp_group_load(&grp, MBEDTLS_ECP_DP_SECP256K1);

  ctx.grp = grp;
  ctx.d = priv;

  size_t siglen = 73;
  unsigned char *sig = new unsigned char[siglen];

  mbedtls_ecdsa_write_signature(&ctx, MBEDTLS_MD_SHA256, hash_data->blob,
                                hash_data->size, sig, &siglen,
                                mbedtls_entropy_func, &entropy);
  // mbedtls_mpi r, s;
  // mbedtls_mpi_init(&r);
  // mbedtls_mpi_init(&s);

  // if (mbedtls_ecdsa_sign(&grp, &r, &s, &priv, hash_data->blob,
  // hash_data->size,
  //                       mbedtls_entropy_func, &entropy) == 0) {
  //  printf("ENCLAVE: Signature generation successful\n");
  //} else {
  //  printf("ENCLAVE: Something went wrong during signing\n");
  //}

  // sig_data->blob = (unsigned char *)std::malloc(65);
  // sig_data->size = 65;

  // std::memset(sig_data->blob, 0, 65);

  // std::memcpy(sig_data->blob, r.p, 32);
  // std::memcpy(sig_data->blob + 32, s.p, 32);

  // printf("Array: \n");
  // for(int i = 0; i < 65; ++i) {
  //        printf("%hu", sig_data[i]);
  //}
  // printf("\n");

  //printf("ENCLAVE: R length: %lu, S length: %lu\n", sizeof(*r.p), sizeof(*s.p));

  sig_data->blob = sig;
  sig_data->size = siglen;
}

bool enclave_verify_secp256k1_sig(data_t *sealed_pubkey, data_t *sig_data,
                                  data_t *sig_hash) {
  printf("ENCLAVE: enclave_verify_secp256k1()");
  point pkey{0, 0};

  Seal seal_bin{OE_SEAL_POLICY_PRODUCT, PATH_WALLET_PUBKEY};

  enclave_get_pubkey(&pkey, sealed_pubkey);

  mbedtls_ecdsa_context ctx;
  mbedtls_ecp_group grp;
  mbedtls_ecp_point p;
  mbedtls_ecdsa_init(&ctx);
  mbedtls_ecp_group_init(&grp);
  mbedtls_ecp_point_init(&p);
  mbedtls_ecp_group_load(&grp, MBEDTLS_ECP_DP_SECP256K1);
  int ret = 0;
  unsigned char *data = NULL;
  size_t data_size = 0;
  Seal seal{OE_SEAL_POLICY_PRODUCT, PATH_WALLET_PUBKEY};
  if (seal.unseal(sealed_pubkey, &data, &data_size) == 0) {
    if ((ret = mbedtls_ecp_point_read_binary(&grp, &p, data, data_size)) != 0) {
      printf(
          "ENCLAVE: ERROR: mbedtls_ecp_point_read_binary() returned with %d\n",
          ret);
    } else {
      printf("ENCLAVE: x: %lu, y: %lu\n", *p.X.p, *p.Y.p);
      ctx.grp = grp;
      ctx.Q = p;
      if ((ret = mbedtls_ecdsa_read_signature(&ctx, sig_hash->blob,
                                              sig_hash->size, sig_data->blob,
                                              sig_data->size)) == OE_OK) {
        printf("ENCLAVE: Signature OK\n");
      } else {
        printf("ENCLAVE: ERROR: Signature invalid\n");
        printf("\t X: %lu,\tY: %lu\nHash: ", *ctx.Q.X.p, *ctx.Q.Y.p);

        for (int i = 0; i < sig_hash->size; ++i) {
          printf("%u", sig_hash->blob[i]);
        }
        printf("\nSignature: ");
        for (int i = 0; i < sig_data->size; ++i) {
          printf("%u", sig_data->blob[i]);
        }
        printf("\n");
      }
    }
  }
  return ret == 0;
}
