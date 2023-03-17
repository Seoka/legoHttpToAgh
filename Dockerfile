# builder image
FROM golang:1.20.2-alpine3.17 as builder
RUN mkdir /build
ADD *.go /build/
ADD go.* /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o legoHttpToAgh .


# generate clean, final image for end users
FROM scratch

ENV ADGUARD_URL http://adguard/
ENV ADGUARD_USER change
ENV ADGUARD_PASS change

EXPOSE 8080

COPY --from=builder /build/legoHttpToAgh .

# executable
ENTRYPOINT [ "./legoHttpToAgh" ]
