#!/bin/bash

MAJVERSION='0'
MINVERSION='0'
BUILDNUM=`git rev-parse --short HEAD`
VERSUFFIX='rel'

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
builddir='./build/release'

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
# full list: windows linux darwin freebsd netbsd openbsd
OSES=(linux darwin)
ARCHS=(amd64 386)
mkdir -p ${builddir}
rm -rf ${builddir}/*
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
		if [ "$os" == "windows" ]; then
			suffix=".exe"
		fi
		env CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ${builddir}/${name}_${os}_${arch}${suffix} .
		if $UPX; then upx --ultra-brute ${builddir}/${name}_${os}_${arch}${suffix}; fi
		tar -C ${builddir} -zcf ${builddir}/${name}_${os}-${arch}-$FULLVERSION.tar.gz ./${name}_${os}_${arch}${suffix}
		$MD5 ${builddir}/${name}_${os}-${arch}-$FULLVERSION.tar.gz
	done
done

# ARM
# full list: 5 6 7
ARMS=(6 7)
for v in ${ARMS[@]}; do
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$v go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ${builddir}/${name}_arm$v .
done
if $UPX; then upx --ultra-brute ${builddir}/${name}_arm*; fi
tar -C ${builddir} -zcf ${builddir}/${name}_arm-$FULLVERSION.tar.gz $(for v in ${ARMS[@]}; do echo -n "./${name}_arm$v ";done)
$MD5 ${builddir}/${name}_arm-$FULLVERSION.tar.gz

# MIPS # go 1.8+ required
env CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ${builddir}/${name}_mipsle .
env CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o ${builddir}/${name}_mips .

if $UPX; then upx --ultra-brute ${builddir}/${name}_mips**; fi
tar -C ${builddir} -zcf ${builddir}/${name}_mipsle-$FULLVERSION.tar.gz ./${name}_mipsle
tar -C ${builddir} -zcf ${builddir}/${name}_mips-$FULLVERSION.tar.gz ./${name}_mips
$MD5 ${builddir}/${name}_mipsle-$FULLVERSION.tar.gz
