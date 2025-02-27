#! /bin/bash

PKGS=$(go list github.com/lianxiangcloud/linkchain/... | grep -v /vendor/ | grep -v /libs/p2p | grep -v /contract| grep -v /evidence | grep -v /wallet/daemon | grep -v /wallet/rpc )

set -e

echo "mode: atomic" > coverage.txt
for pkg in ${PKGS[@]}; do
	go test -timeout 5m -race -coverprofile=profile.out -covermode=atomic "$pkg"
	if [ -f profile.out ]; then
		tail -n +2 profile.out >> coverage.txt;
		rm profile.out
	fi
done
