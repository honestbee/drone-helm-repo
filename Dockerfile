FROM alpine:3.6
MAINTAINER Vincent De Smet <vincent.drl@gmail.com>

ENV VERSION v2.5.1
ENV FILENAME helm-${VERSION}-linux-amd64.tar.gz

RUN set -ex \
  && apk -v --no-cache add \
    curl \
  && curl -sLo /tmp/${FILENAME} http://storage.googleapis.com/kubernetes-helm/${FILENAME} \
  && tar -zxvf /tmp/${FILENAME} -C /tmp \
  && mv /tmp/linux-amd64/helm /bin/helm \
  && rm -rf /tmp 

LABEL description="Helm."
LABEL base="alpine"
LABEL language="golang"


COPY bin/drone-helm-repo /bin/drone-helm-repo

ENTRYPOINT [ "/bin/drone-helm-repo" ]
