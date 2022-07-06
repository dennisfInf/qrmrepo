## Build
The CMakeLists.txt does everything for you.
1. It generates a keypair
2. Signs the enclave
3. Builds the Go-Rest-Server

and more.
```bash
mkdir -p build && cd build && rm -rf * && cmake .. && make

```

## Run
Make sure that you are in the build folder.
```bash
./server
```
