#!/bin/bash

# K8S_RELEASE must be set in your environment, for example:1.17
if [ -n "${K8S_RELEASE}" ]; then
	echo "Setting up reference directories for ${K8S_RELEASE} release."
else
	echo "You must set K8S_RELEASE to a release version, such as 1.17"
	exit 1
fi

ROOTDIR=$(pwd)
echo base dir ${ROOTDIR}

# change <major>.<minor> to <major>_<minor>
echo "${K8S_RELEASE}" | sed "s/\./_/g" > k.tmp

# example: 1_17
VERSION_DIR="$( cat k.tmp )"
echo version ${VERSION_DIR}

echo ${VERSION_DIR} | sed "s/[0-9]_//g" > minor.tmp
echo ${VERSION_DIR} | sed "s/_[0-9]*//g" > major.tmp

MINOR_VERSION="$( cat minor.tmp )"
MAJOR_VERSION="$( cat major.tmp )"

ONE=1
PREV_VERSION_DIR="v""${MAJOR_VERSION}_""$((${MINOR_VERSION} - ${ONE}))"
echo previous version ${PREV_VERSION_DIR}
VERSION_DIR="v${VERSION_DIR}"
echo version ${VERSION_DIR}

# Set up versioned directories for new release

# kubectl-command static-includes, toc.yaml
mkdir -p ./gen-kubectldocs/generators/${VERSION_DIR}

if ! [ -f "${ROOTDIR}/gen-kubectldocs/generators/${VERSION_DIR}/config.yaml" ]; then
		cp -r ${ROOTDIR}/gen-kubectldocs/generators/${PREV_VERSION_DIR}/* ${ROOTDIR}/gen-kubectldocs/generators/${VERSION_DIR}/
fi

# api reference config.yaml, swagger.json
mkdir -p ${ROOTDIR}/gen-apidocs/config/${VERSION_DIR}

# config.yaml
if ! [ -f "${ROOTDIR}/gen-apidocs/config/${VERSION_DIR}/config.yaml" ]; then
	if [ -f "${ROOTDIR}/gen-apidocs/config/${PREV_MINOR_VERSION}/config.yaml" ]; then
			cp ${ROOTDIR}/gen-apidocs/config/${PREV_MINOR_VERSION}/config.yaml ${ROOTDIR}/gen-apidocs/config/${VERSION_DIR}/config.yaml
			echo "Using config file: ${ROOTDIR}/gen-apidocs/config/${PREV_MINOR_VERSION}/config.yaml"
	else
			cp ${ROOTDIR}/gen-apidocs/config/config.yaml ${ROOTDIR}/gen-apidocs/config/${VERSION_DIR}/config.yaml
	fi
fi

# copy versioned files to the base config dir
cp ${ROOTDIR}/gen-apidocs/config/${VERSION_DIR}/config.yaml ${ROOTDIR}/gen-apidocs/config/config.yaml

# swagger.json
# copy swagger file using makefile target

# clean up
rm k.tmp; rm minor.tmp; rm major.tmp