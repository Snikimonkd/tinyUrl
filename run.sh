#!/bin/bash
while getopts ":n" OPTION; do
  case "$OPTION" in
    n)
      echo "in mempory version of programm"
      sudo docker build -t tinyurl .
      sudo docker run -d -p 5000:5000 tinyurl
      exit 0
      ;;
  esac
done

echo "sql version of programm"
sudo docker-compose build
sudo docker-compose up -d