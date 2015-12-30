#!/bin/bash -e
#
#  Cut a new release
#
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do # resolve $SOURCE until the file is no longer a symlink
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE" # if $SOURCE was a relative symlink, we need to resolve it relative to the path where the symlink file was located
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
#
#

if [ -z "${1}" ]; then
  echo "please specify release number"
  read release
else
  release="${1}"
fi

release_name="dozy-${release}_linux-x86-64"
release_folder="${DIR}/${release_name}"

mkdir "${release_folder}"
# build the executable
"${DIR}/build-linux.sh"
# copy the executable and readme
cp "${DIR}/"{dozy,README.md} "${release_folder}"
# pack the release
cd "${DIR}"
tar czf "${release_name}.tar.gz" "${release_name}"
# clean up
rm -r "${release_folder}"

echo "release \"${release}\": ${release_folder}.tar.gz"
