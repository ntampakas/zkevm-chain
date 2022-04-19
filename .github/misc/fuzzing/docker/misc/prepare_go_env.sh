#!/bin/bash
set -e

if [ ! -e /TestCode/.done ]; then

    cd /TestCode
    go mod init zkevmchaintest

    go get -u github.com/ethereum/go-ethereum
    cd /go/pkg/mod/github.com/ethereum/go-ethereum\@v1.10.17/
    make devtools

    cd /TestCode/l1bridge/
    abigen --abi ZkEvmL1Bridge.abi --pkg l1bridge --type ZkEvmL1Bridge --out l1bridge.go
    rm ZkEvmL1Bridge.abi

    cd /TestCode/

cat <<EOF >> go.mod
            require (
                l1bridge v1.0.0
            )

            replace (
                l1bridge v1.0.0 => ./l1bridge
            )
EOF

    cd l1bridge
    go mod init l1bridge
    go mod tidy

    cd ..
    go mod tidy

    rm l1bridge/go.*

    touch /TestCode/.done
fi
