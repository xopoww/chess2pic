# syntax=docker/dockerfile:1

FROM alpine/git as clone

WORKDIR /repo
RUN git clone https://github.com/xopoww/chess2pic.git .


FROM golang:1.17-alpine as build

WORKDIR /chess2pic
COPY --from=clone /repo/ ./
RUN go build -o ./build/chess2pic-api-server ./cmd/chess2pic-api-server

FROM alpine:latest

WORKDIR /app
COPY --from=build /chess2pic/build/chess2pic-api-server .

EXPOSE 8080
ENV HOST=0.0.0.0 PORT=8080
CMD [ "./chess2pic-api-server" ]