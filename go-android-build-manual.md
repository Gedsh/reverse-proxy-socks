#Build for android

ndk_version=23.1.7779620

NDK="$HOME/Android/Sdk/ndk/${ndk_version}"
export PATH="$PATH:$NDK/toolchains/llvm/prebuilt/linux-x86_64/bin"
export GOPATH=$HOME/go/
export GOOS='android'
export CGO_ENABLED=1

#armv7a:
export GOARCH='arm'
export GOARM='7'
export CC="armv7a-linux-androideabi16-clang"
export CCX="armv7a-linux-androideabi16-clang++"
export CGO_CFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_CPPFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_CXXFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_FFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_LDFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"

#arm64:
export GOARCH='arm64'
export GOARM64='v8.0,crypto'
export CC="aarch64-linux-android21-clang"
export CCX="aarch64-linux-android21-clang++"
export CGO_CFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_CPPFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_CXXFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_FFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_LDFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"

#x86_64:
export GOARCH='amd64'
export GOAMD64='v2'
export CC="x86_64-linux-android21-clang"
export CCX="x86_64-linux-android21-clang++"
export CGO_CFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_CPPFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_CXXFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_FFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"
export CGO_LDFLAGS="-g -O3 -funsafe-math-optimizations -ftree-vectorize -fvectorize -fslp-vectorize"

#common:

go clean
go build -x -ldflags="-s -w" -compiler gc -gcflags="-m -dwarf=false" -o libreverseproxy.so

