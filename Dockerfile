FROM ubuntu:latest

RUN apt-get update \
      && apt-get install -y --no-install-recommends \
         make \
         sudo \
         gcc \
         make \
         binutils \
         libc6-dev \
         gdb \
         wget \
      && apt-get clean -y \
      && rm -rf /var/lib/apt/lists/*

RUN wget --no-check-certificate https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz \
      && sudo tar -C /usr/local -xzf go1.11.5.linux-amd64.tar.gz \
      && rm -rf go1.11.5.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin
