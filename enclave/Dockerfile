FROM ubuntu:20.04
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /root
RUN apt update && apt install -y cmake \
	make \
	wget \
	git \
	neovim \
	gnupg2 \
	clang

# Install go
RUN wget -q -O- https://go.dev/dl/go1.18.3.linux-amd64.tar.gz | tar -xz -C /usr/local
ENV PATH="${PATH}:/usr/local/go/bin"

# Install neovim
RUN apt install -y ninja-build gettext libtool libtool-bin autoconf automake cmake g++ pkg-config unzip curl doxygen
RUN git clone https://github.com/neovim/neovim && cd neovim && make install && cd .. && rm -rf neovim
ENV PATH="${PATH}:/usr/local/bin"

# Install openenclave
RUN echo 'deb [arch=amd64] https://download.01.org/intel-sgx/sgx_repo/ubuntu focal main' | tee /etc/apt/sources.list.d/intel-sgx.list && \
	wget -qO - https://download.01.org/intel-sgx/sgx_repo/ubuntu/intel-sgx-deb.key | apt-key add - && \
	echo "deb http://apt.llvm.org/focal/ llvm-toolchain-focal-10 main" | tee /etc/apt/sources.list.d/llvm-toolchain-focal-10.list && \
	wget -qO - https://apt.llvm.org/llvm-snapshot.gpg.key | apt-key add - && \
	echo "deb [arch=amd64] https://packages.microsoft.com/ubuntu/20.04/prod focal main" | tee /etc/apt/sources.list.d/msprod.list && \
	wget -qO - https://packages.microsoft.com/keys/microsoft.asc |apt-key add - && apt update

RUN echo "source /opt/openenclave/share/openenclave/openenclaverc" >> /root/.bashrc

RUN apt -y install clang-10 \
	libssl-dev \
	gdb \
	libsgx-enclave-common \
	libsgx-quote-ex \
	libprotobuf17 \
	libsgx-dcap-ql \
	libsgx-dcap-ql-dev \
	az-dcap-client \
	open-enclave

RUN apt install -y clang-format

WORKDIR /server

RUN go install golang.org/x/tools/gopls@latest

RUN mkdir -p build && cd build && rm -rf * && cmake .. && make

RUN cd build && ./server

ENTRYPOINT [ "bash" ]

