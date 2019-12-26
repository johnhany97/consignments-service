FROM alpine:latest

RUN mkdir /app
WORKDIR /app
ADD consignments-service /app/consignments-service

CMD ["./consignments-service"]