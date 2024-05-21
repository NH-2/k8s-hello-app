FROM --platform=amd64 golang:alpine3.19 AS builder

# Create appuser.
ENV USER=appuser
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

COPY . /app
WORKDIR /app
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o hello .

FROM scratch AS runner

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app/hello /hello

# Use an unprivileged user.
USER appuser:appuser

CMD [ "/hello" ]
