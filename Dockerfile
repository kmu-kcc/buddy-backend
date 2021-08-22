FROM golang:1.17
WORKDIR /bin/
COPY ./buddy ./
RUN chmod +x ./buddy
EXPOSE 3000
CMD ["./buddy", "--port", "3000"]
