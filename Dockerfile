FROM golang:1.23-alpine

WORKDIR /app

COPY . .

ENV CGO_ENABLED=0
RUN go build -o /bin/ldapb ./cmd

FROM scratch

COPY --from=0 /bin/ldapb /bin/ldapb

EXPOSE 389

ENTRYPOINT ["ldapb"]
