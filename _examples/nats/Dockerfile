FROM alpine:3.18.4 as rtl_build
RUN apk --no-cache add git cmake libusb-dev make gcc g++ alpine-sdk
WORKDIR /tmp
RUN git clone https://github.com/minghsu0107/librtlsdr
WORKDIR /tmp/librtlsdr

RUN mkdir build && cd build && cmake ../ && make && make install
RUN ls /usr/local/bin/rtl_*

FROM golang:1.20-alpine as go_build
RUN mkdir -p /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app .

FROM alpine:3.18.4
RUN apk --no-cache add alsa-utils libusb ca-certificates curl jq bash
COPY --from=go_build /bin/app /bin/app
COPY --from=rtl_build /usr/local/bin/rtl_test /bin/rtl_test
COPY --from=rtl_build /usr/local/bin/rtl_fm /bin/rtl_fm
COPY --from=rtl_build /usr/local/bin/rtl_power /bin/rtl_power
COPY --from=rtl_build /usr/local/lib/librtlsdr.so.0 /usr/lib/librtlsdr.so.0
WORKDIR /
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
ENTRYPOINT ["/bin/app"]
