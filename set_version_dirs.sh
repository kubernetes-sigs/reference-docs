#!/bin/bash

# K8S_RELEASE must be set in your environment, for example, 1.17
if [ -n "${K8S_RELEASE}" ]; then
	echo "Setting up reference directories for ${K8S_RELEASE} release."
else
	echo "You must set K8S_RELEASE to a release version, such as 1.17"
	exit 1
fi

ROOTDIR=$(pwd)
echo base dir ${ROOTDIR}

# change <major>.<minor> to <major>_<minor>
VERSION_DIR="$(echo "${K8S_RELEASE}" | sed "s/\./_/g")"
echo version ${VERSION_DIR}

MINOR_VERSION="$(echo ${VERSION_DIR} | sed "s/[0-9]_//g")"
echo minor version ${MINOR_VERSION}

MAJOR_VERSION="$(echo ${VERSION_DIR} | sed "s/_[0-9]*//g")"
echo major version ${MAJOR_VERSION}

ONE=1

PREV_VERSION_DIR="v""${MAJOR_VERSION}_""$((${MINOR_VERSION} - ${ONE}))"
echo previous version dir ${PREV_VERSION_DIR}

VERSION_DIR="v${VERSION_DIR}"
echo version dir ${VERSION_DIR}

# Set up versioned directories for new release

# kubectl-command static-includes, toc.yaml
mkdir -p ./gen-kubectldocs/generators/${VERSION_DIR}

if ! [ -f "${ROOTDIR}/gen-kubectldocs/generators/${VERSION_DIR}/toc.yaml" ]; then
		cp -r ${ROOTDIR}/gen-kubectldocs/generators/${PREV_VERSION_DIR}/* ${ROOTDIR}/gen-kubectldocs/generators/${VERSION_DIR}/
fi

# api reference config.yaml
mkdir -p ${ROOTDIR}/gen-apidocs/config/${VERSION_DIR}

# config.yaml
if ! [ -f "${ROOTDIR}/gen-apidocs/config/${VERSION_DIR}/config.yaml" ]; then
	if [ -f "${ROOTDIR}/gen-apidocs/config/${PREV_VERSION_DIR}/config.yaml" ]; then
			cp ${ROOTDIR}/gen-apidocs/config/${PREV_VERSION_DIR}/config.yaml ${ROOTDIR}/gen-apidocs/config/${VERSION_DIR}/config.yaml
			echo "Using config file: ${ROOTDIR}/gen-apidocs/config/${PREV_VERSION_DIR}/config.yaml"
	else
			cp ${ROOTDIR}/gen-apidocs/config/config.yaml ${ROOTDIR}/gen-apidocs/config/${VERSION_DIR}/config.yaml
	fi
fi

# copy versioned files to the base config dir
# revisit
cp ${ROOTDIR}/gen-apidocs/config/${VERSION_DIR}/config.yaml ${ROOTDIR}/gen-apidocs/config/config.yaml

# api reference swagger.json is copied by updateapispec make target
