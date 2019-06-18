FROM circleci/golang:1.12

ENV PATH "$PATH:/home/circleci/.local/bin"

USER root
RUN apt-get update \
 && apt-get install -y python-pip \
 && pip install awscli

RUN mkdir -p /tmp/app

COPY ./ /tmp/app/

RUN cd /tmp/app/ \
        && go mod vendor \
        && go build \
        && cp ops-cli /usr/local/bin/ \
        && cd / \
        && rm -rf /tmp/app
USER circleci
