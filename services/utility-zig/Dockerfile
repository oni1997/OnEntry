FROM alpine:3.19 AS builder
RUN apk add --no-cache zig build-base
WORKDIR /build
COPY build.zig .
COPY src ./src
RUN zig build -Doptimize=ReleaseSmall

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /build/zig-out/bin/utility-service .
EXPOSE 8084
CMD ["./utility-service"]
