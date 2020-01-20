FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY server /
COPY assets /assets
EXPOSE 8080
CMD ["/server"]