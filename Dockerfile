# Multi-stage build for SpiralDB with TRUE FHE
FROM ubuntu:22.04 AS builder

# Install dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    golang-go \
    cmake \
    git \
    && rm -rf /var/lib/apt/lists/*

# Build SEAL
WORKDIR /build
RUN git clone https://github.com/microsoft/SEAL.git && \
    cd SEAL && \
    git checkout 4.1.2 && \
    cmake -S . -B build -DSEAL_BUILD_SEAL_C=ON -DCMAKE_INSTALL_PREFIX=/usr/local && \
    cmake --build build && \
    cmake --install build

# Copy SpiralDB source
WORKDIR /app
COPY . .

# Build FHE bridge
RUN g++ -std=c++17 -O3 -fPIC -shared \
    -I/usr/local/include/SEAL-4.1 \
    -o libfhe_bridge.so \
    fhe_bridge.cpp \
    /usr/local/lib/libseal-4.1.a

# Build Go app
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-I. -I/usr/local/include/SEAL-4.1"
ENV CGO_LDFLAGS="-L. -lfhe_bridge -L/usr/local/lib -lseal-4.1 -lstdc++ -lm"
RUN go build -o spiraldb .

# Final stage
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y \
    libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/spiraldb .
COPY --from=builder /app/libfhe_bridge.so .
COPY --from=builder /usr/local/lib/libseal-4.1.a /usr/local/lib/

ENV LD_LIBRARY_PATH=/app:/usr/local/lib
EXPOSE 8094
CMD ["./spiraldb"]
