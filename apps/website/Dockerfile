FROM node:23-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .

ARG VITE_DOMAIN
ARG VITE_API_URL
ARG VITE_BOT_NAME
ENV VITE_DOMAIN=$VITE_DOMAIN \
    VITE_API_URL=$VITE_API_URL \
    VITE_BOT_NAME=$VITE_BOT_NAME

RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
