FROM ubuntu:22.04

RUN apt-get update -y  && \
  apt-get install -y \
  mysql-client \
  golang-go \
  && apt-get clean

WORKDIR /app

COPY . .

RUN chmod +x script.sh apache.sh

EXPOSE 8000

CMD ["bash", "-c", "./script.sh && ./apache.sh"]