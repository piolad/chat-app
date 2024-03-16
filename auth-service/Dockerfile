FROM rust

RUN mkdir -p /home/app
COPY . /home/app

RUN apt-get update && \
    apt-get install -y protobuf-compiler


CMD [ "sh", "-c", "cd /home/app && cargo run" ]
