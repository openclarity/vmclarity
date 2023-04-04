FROM alpine:3.17.0

ARG IMAGE_VERSION

LABEL name='goat-image'
LABEL vendor='Openclarity'
LABEL version=${IMAGE_VERSION}
LABEL release=${IMAGE_VERSION}

# add GNU ps for -C, -o cmd, and --no-headers support
# https://github.com/aquasecurity/kube-bench/issues/109
RUN apk --no-cache add procps

# Upgrading apk-tools to remediate CVE-2021-36159 - https://snyk.io/vuln/SNYK-ALPINE314-APKTOOLS-1533752
# https://github.com/aquasecurity/kube-bench/issues/943
RUN apk --no-cache upgrade apk-tools

COPY goat.sh /app/goat.sh
COPY families /app/families
RUN chmod +x /app/goat.sh
WORKDIR /app

ENTRYPOINT ["/app/goat.sh"]
