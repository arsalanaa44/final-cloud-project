FROM golang As builder 

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM alpine

CMD ["/root/server"]

WORKDIR /root/

COPY --from=builder /app/server . 