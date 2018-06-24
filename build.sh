#!/bin/bash

version_gt() {
    test "$(printf '%s\n' "$@" | sort -V | head -n 1)" != "$1";
}

scriptpath="$(cd "$(dirname "$0")"; pwd -P)"
go=$(command -v go)

command -v go >/dev/null 2>&1 || { echo >&2 "Go 1.8+ needs to be available for compilation."; exit 1; }
command -v git >/dev/null 2>&1 || { echo >&2 "Git needs to be available for compilation."; exit 1; }

bos=$($go run $scriptpath/build/build.go -o)
barch=$($go run $scriptpath/build/build.go -a)

os="$bos"
arch="$barch"
cmd="-buildmode=exe"

out="grm"

i=0
params=("$@")
for arg in "$@"; do
    case "$arg" in
        --help|-h )
            echo "Usage $0 [ posix or GNU style options ]"
            echo -e "-f|--force\t\t\t\tForce recompilation of all packages"
            echo -e "-a|--arch\033[3m[=]target_arch\033[0m\t\t\t\tSelect target architecture (amd64, arm)"
            echo -e "-o|--os\033[3m[=]target_os\033[0m\t\t\t\tSelect the target operating system (linux, darwin, windows, freebsd)"
            echo -e "-v|--verbose\t\t\t\tEnable verbose compilation mode"
            echo -e "--ogo\033[3m[=]path_to go_binary\033[0m\t\t\t\tSelect a different Go binary for compilation"
            exit 0
            ;;
        --force|-f )
            cmd="-a $cmd"
            ;;
        --verbose|-v )
            cmd="-v $cmd"
            ;;
        --arch=*|-a=* )
            arch=`echo $arg | sed 's/[-a-zA-Z0-9]*=//'`
            ;;
        --arch|-a )
            ((i++))
            arch="${params[$i]}"
            shift
            ;;
        --os=* )
            os=`echo $arg | sed 's/[-a-zA-Z0-9]*=//'`
            ;;
        --os )
            ((i++))
            os="${params[$i]}"
            shift
            ;;
        --go=* )
            go=`echo $arg | sed 's/[-a-zA-Z0-9]*=//'`
            ;;
        --go )
            ((i++))
            go="${params[$i]}"
            shift
            ;;
      esac
      ((i++))
done

bversion=$($go run $scriptpath/build/build.go -v)

if version_gt "1.8.0" ${bversion}; then
    echo "Go version 1.8 or later is required. Found Go version: $bversion"
    exit 1
fi

if [[ "$os" == "windows" ]]; then
    out="$out.exe"
fi

gitrev=$(git rev-parse HEAD)
buildversion=$(git describe --tags >/dev/null 2>1||echo ${gitrev})
builddate=$(date "+%Y-%m-%dT%H:%M:%S+%Z")
ldflags="-X=main.buildVersion=$buildversion -X=main.buildDate=$builddate"

echo "########################################################################################"
echo "# Build OS: $bos"
echo "# Build Arch: $barch"
echo "# Go Version: $bversion"
echo "# Target OS: $os"
echo "# Target Arch: $arch"
echo "# BuildTime: $builddate"
echo "# Build Version: $buildversion"
echo "########################################################################################"

mkdir -p target/${os}
export GOPATH="$scriptpath"
${go} build -ldflags="$ldflags" ${cmd} -o target/$os/$out grm
chmod +x target/${os}/grm