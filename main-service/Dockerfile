FROM node

RUN mkdir -p /home/app
COPY . /home/app


WORKDIR /home/app/src

RUN npm install
CMD "npm" "start"