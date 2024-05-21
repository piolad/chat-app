FROM rust

RUN rustup component add rustfmt

RUN mkdir -p /home/app 
COPY . /home/app

WORKDIR /home/app

CMD ["cargo", "run"]
