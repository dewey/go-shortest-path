FROM alpine:latest
RUN apk add --update ca-certificates
COPY api api
EXPOSE 8080
CMD ["./api"]