FROM circleci/golang:1.12

ENV PATH "$PATH:/home/circleci/.local/bin"

USER root
RUN apt-get update \
 && apt-get install -y python-pip curl \
 && pip install awscli

RUN curl -Os https://releases.hashicorp.com/terraform/0.12.29/terraform_0.12.29_linux_amd64.zip \
	&& unzip terraform_0.12.29_linux_amd64.zip -d /usr/local/bin \
	&& rm terraform_0.12.29_linux_amd64.zip

RUN mkdir -p /tmp/app

COPY ./ /tmp/app/

RUN cd /tmp/app/ \
        && go mod vendor \
        && go build \
        && cp ops-cli /usr/local/bin/ \
        && cd / \
        && rm -rf /tmp/app
USER circleci
