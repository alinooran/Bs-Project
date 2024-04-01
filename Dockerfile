FROM golang:1.20.5-alpine3.18
WORKDIR /app
COPY go* .
RUN go mod download
COPY . .
ENV DB_HOST="host.docker.internal"
ENV DB_USER="root"
ENV DB_PASSWORD="117101246950"
ENV DB_NAME="bs_project"
ENV DB_PORT="3306"
ENV JWT_SECRET_KEY="secret"
ENV HASH_COST=15
ENV EMAIL="sbuproject403@gmail.com"
ENV EMAIL_PASS="ekvljsbtzjmparoy"
CMD ["go", "run", "main.go"]