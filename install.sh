#!/bin/sh

set -e

downloadFolder="${HOME}/Downloads"
owner="axetroy"
repo="dvm"
exe_name="dvm"

mkdir -p ${downloadFolder}

get_arch() {
    a=$(uname -m)
    case ${a} in
    "x86_64" | "amd64" )
        echo "amd64"
        ;;
    "i386" | "i486" | "i586")
        echo "386"
        ;;
    "aarch64" | "arm64" | "arm")
        echo "arm64"
        ;;
    *)
        echo ${NIL}
        ;;
    esac
}

get_os(){
    echo $(uname -s | awk '{print tolower($0)}')
}

os=$(get_os)
arch=$(get_arch)
dest_file="${downloadFolder}/${exe_name}_${os}_${arch}.tar.gz"

if [ $# -eq 0 ]; then
    asset_path=$(
        command curl -sSf https://github.com/${owner}/${repo}/releases |
            command grep -o "/${owner}/${repo}/releases/download/.*/${exe_name}_${os}_${arch}.tar.gz" |
            command head -n 1
    )
    if [[ ! "$asset_path" ]]; then exit 1; fi
    asset_uri="https://github.com${asset_path}"
else
    asset_uri="https://github.com/${owner}/${repo}/releases/download/${1}/${exe_name}_${os}_${arch}.tar.gz"
fi

mkdir -p ${downloadFolder}

echo "[1/3] Download ${asset_uri} to ${downloadFolder}"
rm -f ${dest_file}
curl --fail --location --output "${dest_file}" "${asset_uri}"

binDir=/usr/local/bin

echo "[2/3] Install ${exe_name} to the ${binDir}"
mkdir -p ${HOME}/bin
tar -xz -f ${dest_file} -C ${binDir}
exe=${binDir}/${exe_name}
chmod +x ${exe}

echo "[3/3] Set environment variables"
echo "${exe_name} was installed successfully to ${exe}"
if command -v $exe_name --version >/dev/null; then
    echo "Run '$exe_name --help' to get started"
else
    echo "Manually add the directory to your \$HOME/.bash_profile (or similar)"
    echo "  export PATH=${HOME}/bin:\$PATH"
    echo "Run '$exe --help' to get started"
fi

exit 0