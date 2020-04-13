#!/bin/bash

MAJVERSION='0'
MINVERSION='0'
BUILDNUM=`git rev-parse --short HEAD`
VERSUFFIX='dev'

FULLVERSION="${MAJVERSION}.${MINVERSION}-${VERSUFFIX}-${BUILDNUM}"

while [ $# -gt 0 ]; do
  case "$1" in
    -M|--major=*)
      MAJVERSION="${1#*=}"
      ;;
    -m|--minor=*)
      MINVERSION="${1#*=}"
	  ;;
    *)
      printf "***************************\n"
      printf "* Error: Invalid argument: ${1}: ${1#*=}*\n"
      printf "***************************\n"
      exit 1
  esac
  shift
done

name='cloudflare-ddns'
builddir='./build/dev'

MD5='md5sum'
if [[ "$(uname)" == 'Darwin' ]]; then
	MD5='md5'
fi

UPX=false
if hash upx 2>/dev/null; then
	UPX=true
fi

LDFLAGS="-X main.majVersion=$MAJVERSION -X main.minVersion=$MINVERSION -X main.buildNum=$BUILDNUM -X main.verSuffix=$VERSUFFIX -s -w"
GCFLAGS=""

# X86
OSES=(darwin)
ARCHS=(amd64)
mkdir -p ${builddir}
rm -rf ${builddir}/*
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
		if [ "$os" == "windows" ]; then
			suffix=".exe"
			LDFLAGS="-X main.majVersion=$MAJVERSION -X main.minVersion=$MINVERSION -X main.buildNum=$BUILDNUM -X main.verSuffix=$VERSUFFIX -s -w"
		fi
		env CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ${builddir}/${name}_${os}_${arch}${suffix} .
		if $UPX; then upx -9 ${builddir}/${name}_${os}_${arch}${suffix}; fi
	done
done
