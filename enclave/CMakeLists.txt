cmake_minimum_required(VERSION 3.13)
project("Elon_Mask_Enclave")
find_package(OpenEnclave CONFIG REQUIRED)

set(CMAKE_C_COMPILER clang)
set(CMAKE_CXX_COMPILER clang++)
set(CMAKE_EXPORT_COMPILE_COMMANDS true)
set(CMAKE_CXX_STANDARD 11)
set(OE_CRYPTO_LIB mbedtls CACHE STRING "Crypto library used by enclaves.")

add_subdirectory(enclave)
add_subdirectory(host)

# Generate key
add_custom_command(
  OUTPUT private.pem public.pem
  COMMAND openssl genrsa -out private.pem -3 3072
  COMMAND openssl rsa -in private.pem -pubout -out public.pem)

# Sign enclave
add_custom_command(
  OUTPUT enclave/enclave.signed
  DEPENDS enclave enclave/project.conf private.pem
  COMMENT "sign enclave"
  COMMAND openenclave::oesign sign -e $<TARGET_FILE:enclave> -c ${CMAKE_SOURCE_DIR}/enclave/project.conf -k private.pem)

add_custom_command(
	OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/server
	DEPENDS host
	COMMENT "build webauthn server"
	WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}/api
	COMMAND go build -o ${CMAKE_CURRENT_BINARY_DIR}/server ${CMAKE_SOURCE_DIR}/api)

add_custom_target(sign ALL DEPENDS enclave/enclave.signed)
add_custom_target(webauthn_api ALL DEPENDS ${CMAKE_CURRENT_BINARY_DIR}/server)
