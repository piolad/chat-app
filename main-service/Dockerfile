FROM node

RUN mkdir -p /home/app
COPY . /home/app


WORKDIR /home/app/src

# CMD ["tail", "-f", "/dev/null"]
RUN npm install
CMD "npm" "start"