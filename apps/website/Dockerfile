FROM node:23-alpine AS build
WORKDIR /app

RUN corepack enable

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .

ARG VITE_DOMAIN
ARG VITE_API_URL
ARG VITE_BOT_NAME
ENV VITE_DOMAIN=$VITE_DOMAIN \
    VITE_API_URL=$VITE_API_URL \
    VITE_BOT_NAME=$VITE_BOT_NAME

RUN pnpm run build

FROM nginx:1.29.1-alpine-slim
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
