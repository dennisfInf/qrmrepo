add_custom_command(
  OUTPUT project_u.h project_u.c project_args.h
  DEPENDS ${CMAKE_SOURCE_DIR}/project.edl
  COMMAND openenclave::oeedger8r --untrusted ${CMAKE_SOURCE_DIR}/project.edl --search-path ${OE_INCLUDEDIR} --search-path ${OE_INCLUDEDIR}/openenclave/edl/sgx)

add_library(host host.cpp host.h ${CMAKE_CURRENT_BINARY_DIR}/project_u.c)
target_compile_options(host PUBLIC -fpie)

# Needed by the generated MyEnclave_u.c
target_include_directories(host PRIVATE ${CMAKE_CURRENT_BINARY_DIR})

# Include generated *.pb.h files
include_directories("${CMAKE_CURRENT_BINARY_DIR}")

target_link_libraries(host openenclave::oehost)
