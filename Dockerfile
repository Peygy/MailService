# Back go
FROM golang:1.23.6-alpine AS go-builder
WORKDIR /app
COPY ./back/go.mod ./back/go.sum ./
RUN go mod download
COPY back/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server

# Client react
FROM node:18 AS react-builder
WORKDIR /app
COPY front/package.json front/package-lock.json ./
RUN npm install
COPY front/ .
RUN npm run build

# Nginx
FROM nginx:alpine
WORKDIR /app
COPY --from=go-builder /app/server /app/server
COPY --from=react-builder /app/build /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY start.sh /app/
RUN chmod +x /app/start.sh
EXPOSE 80
CMD ["/app/start.sh"]