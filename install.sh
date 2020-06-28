#!/usr/bin/env bash
set -e

downloadFolder="${HOME}/Downloads"

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
    *)
        echo ${NIL}
        ;;
    esac
}

get_os(){
    echo $(uname -s | awk '{print tolower($0)}')
}

main() {
    local os=$(get_os)
    local arch=$(get_arch)
    local dest_file="${downloadFolder}/dvm_${os}_${arch}.tar.gz"

    if [[ $# -eq 0 ]]; then
        asset_path=$(
            command curl -sSf https://github.com/axetroy/dvm/releases |
                command grep -o "/axetroy/dvm/releases/download/.*/dvm_${os}_${arch}\\.tar.gz" |
                command head -n 1
        )
        if [[ ! "$asset_path" ]]; then exit 1; fi
        asset_uri="https://github.com${asset_path}"
    else
        asset_uri="https://github.com/axetroy/dvm/releases/download/${1}/dvm_${os}_${arch}\\.tar.gz"
    fi

    mkdir -p ${downloadFolder}

    echo "[1/3] Download ${asset_uri} to ${downloadFolder}"
    rm -f ${dest_file}
    # wget -P "${downloadFolder}" "${asset_uri}"
    curl --location --output "${dest_file}" "${asset_uri}"

    binDir=/usr/local/bin

    echo "[2/3] Install dvm to the ${binDir}"
    mkdir -p ${HOME}/bin
    tar -xz -f ${dest_file} -C ${binDir}
    exe=${binDir}/dvm
    chmod +x ${exe}

    echo "[3/3] Set environment variables"
    echo "dvm was installed successfully to ${exe}"
    if command -v dvm --version >/dev/null; then
        echo "Run 'dvm --help' to get started"
    else
        echo "Manually add the directory to your \$HOME/.bash_profile (or similar)"
        echo "  export PATH=${HOME}/bin:\$PATH"
        echo "Run '$exe --help' to get started"
    fi

    exit 0
}

main