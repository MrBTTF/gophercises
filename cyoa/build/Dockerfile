FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY main /
COPY assets /assets
EXPOSE 8080
CMD ["/main"]