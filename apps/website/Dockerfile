FROM node:26-alpine AS build
WORKDIR /app

RUN npm install -g corepack && corepack enable

COPY pnpm-workspace.yaml package.json pnpm-lock.yaml ./
ENV CI=true
RUN pnpm install --frozen-lockfile
COPY . .

ARG VITE_DOMAIN
ARG VITE_API_URL
ARG VITE_BOT_NAME
ARG VITE_TENOR_API_KEY
ENV VITE_DOMAIN=$VITE_DOMAIN \
    VITE_API_URL=$VITE_API_URL \
    VITE_BOT_NAME=$VITE_BOT_NAME \
    VITE_TENOR_API_KEY=$VITE_TENOR_API_KEY

RUN pnpm run build

FROM nginx:1.31.0-alpine-slim
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
