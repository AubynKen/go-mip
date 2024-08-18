# this file is a direct translation of Taskfile.yml for my C++ IDE's intellisense to work better
# if I forget to delete this, just ignore it and use Taskfile

cmake_minimum_required(VERSION 3.10)
project(BridgeProject)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_POSITION_INDEPENDENT_CODE ON)

# Add include directories
include_directories(/usr/local/include)

# Add link directories
link_directories(/usr/local/lib)

# Create the bridge library
add_library(bridge SHARED bridge/bridge.cpp)
target_link_libraries(bridge ortools)

# Set rpath for the shared library
set_target_properties(bridge PROPERTIES
        INSTALL_RPATH "/usr/local/lib"
        BUILD_WITH_INSTALL_RPATH TRUE
)

# Custom target for running the Go program
add_custom_target(run
        COMMAND go run main.go
        DEPENDS bridge
        WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
)

# Custom target for cleaning
add_custom_target(clean_custom
        COMMAND ${CMAKE_COMMAND} -E remove bridge/bridge.o bridge/libbridge.dylib
        WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
)
