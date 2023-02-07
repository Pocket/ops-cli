FROM circleci/golang:1.17

ENV PATH "$PATH:/home/circleci/.local/bin"

RUN mkdir -p /tmp/app

COPY ./ /tmp/app/

USER root
RUN apt-get update \
 && apt-get install -y python3-pip curl \
 && pip install awscli

RUN cd /tmp/app/ \
        && go mod vendor \
        && go build \
        && cp ops-cli /usr/local/bin/ \
        && cd / \
        && rm -rf /tmp/app

# Terraform Setup
USER circleci
RUN git clone https://github.com/tfutils/tfenv.git ~/.tfenv
ENV HOME "/home/circleci"
ENV PATH "$PATH:$HOME/.tfenv/bin"
RUN tfenv install 0.12.29 && tfenv use
