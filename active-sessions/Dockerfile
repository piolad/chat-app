FROM rust

RUN rustup component add rustfmt

RUN mkdir -p /home/app 
COPY . /home/app

RUN apt-get update && \
    apt-get install -y protobuf-compiler

WORKDIR /home/app

CMD ["cargo", "run"]
