FROM alpine:3.19 AS builder
RUN apk add --no-cache build-base wget
RUN wget -q https://ziglang.org/download/0.10.1/zig-linux-x86_64-0.10.1.tar.xz && \
    tar xf zig-linux-x86_64-0.10.1.tar.xz && \
    ln -s /zig-linux-x86_64-0.10.1/zig /usr/local/bin/zig && \
    rm zig-linux-x86_64-0.10.1.tar.xz
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
