FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM alpine:3.21 AS dist

WORKDIR /dist
COPY --from=builder /app/dist/ .

CMD ["sh", "-c", "rm -rf /srv/* && cp -a /dist/. /srv/"]
