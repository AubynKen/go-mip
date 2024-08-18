FROM --platform=linux/amd64 golang:1.22.0

ARG ORTOOLS_VERSION=9.10
ARG ORTOOLS_BUILD=4067

# Combine RUN commands to reduce layers and use `apt-get` best practices
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    lsb-release \
    wget \
    curl \
    && rm -rf /var/lib/apt/lists/* \
    && sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

# Install OR-Tools binaries
RUN wget --no-check-certificate -O or-tools.tar.gz https://github.com/google/or-tools/releases/download/v${ORTOOLS_VERSION}/or-tools_amd64_debian-11_cpp_v${ORTOOLS_VERSION}.${ORTOOLS_BUILD}.tar.gz \
    && tar -xzvf or-tools.tar.gz \
    && mv or-tools_* /opt/or-tools \
    && rm or-tools.tar.gz

# Link OR-Tools libraries
RUN if [ -d "/opt/or-tools/lib" ]; then \
        ln -s /opt/or-tools/lib/* /usr/local/lib/ && \
        ldconfig && \
        echo "OR-Tools libraries linked successfully"; \
    else \
        echo "Error: /opt/or-tools/lib directory not found" && \
        exit 1; \
    fi

WORKDIR /app

# Copy the bridging layer C++ code
COPY bridge /app/bridge

# Set environment variables
ENV CGO_CXXFLAGS="-I/opt/or-tools/include -I/opt/or-tools/lib"
ENV CGO_LDFLAGS="-L/opt/or-tools/lib -L/app/bridge -Wl,-rpath,/opt/or-tools/lib -Wl,-rpath,/app/bridge"
ENV LD_LIBRARY_PATH="/opt/or-tools/lib:/app/bridge:$LD_LIBRARY_PATH"

# Build the bridging layer in C++
CMD task clean && task build-bridge

# Copy the rest of the source code
COPY . .

# Build the Go application and run it
CMD task run
