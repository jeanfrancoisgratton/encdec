#!/usr/bin/env bash

PKGDIR=encdec-1.31.00-0_amd64

mkdir -p ${PKGDIR}/opt/bin ${PKGDIR}/DEBIAN
mv control ${PKGDIR}/DEBIAN/
mv preinst ${PKGDIR}/DEBIAN/

echo "Building binary from source"
cd ../src
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -buildid=" -o ../__debian/${PKGDIR}/opt/bin/encdec .
sudo chown 0:0 ../__debian/${PKGDIR}/opt/bin/encdec

echo "Binary built. Now packaging..."
cd ../__debian/
dpkg-deb -b ${PKGDIR}
