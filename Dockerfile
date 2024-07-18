FROM golang:1.20.5-alpine3.18
RUN addgroup app && adduser -S -G app app
USER app
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
EXPOSE 8080
ENV DB_HOST="host.docker.internal"
ENV DB_USER="root"
ENV DB_PASSWORD="117101246950"
ENV DB_NAME="bs_project"
ENV DB_PORT="3306"
ENV JWT_SECRET_KEY="jwtsecretkey"
ENV HASH_COST=15
CMD ["go", "run", "main.go"]