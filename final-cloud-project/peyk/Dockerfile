FROM golang As builder 

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o peyk .

FROM alpine

CMD ["/root/peyk"]

WORKDIR /root/

COPY --from=builder /app/peyk . 
