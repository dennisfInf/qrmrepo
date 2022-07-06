add_custom_command(
  OUTPUT project_t.h project_t.c project_args.h
  DEPENDS ${CMAKE_SOURCE_DIR}/project.edl
  COMMAND openenclave::oeedger8r --trusted ${CMAKE_SOURCE_DIR}/project.edl --search-path ${OE_INCLUDEDIR} --search-path ${OE_INCLUDEDIR}/openenclave/edl/sgx)

add_executable(enclave enclave.cpp ${CMAKE_CURRENT_BINARY_DIR}/project_t.c)
# Current API version
target_compile_definitions(enclave PUBLIC OE_API_VERSION=2)

# Needed by the generated project_t.c
target_include_directories(enclave PRIVATE ${CMAKE_CURRENT_BINARY_DIR})

target_link_libraries(enclave openenclave::oeenclave openenclave::oelibcxx openenclave::oecryptombedtls openenclave::oeseal_gcmaes)
