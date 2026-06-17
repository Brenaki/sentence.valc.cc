# syntax=docker/dockerfile:1

# --- build stage ---
FROM node:22-alpine AS build
WORKDIR /app

# VITE_API_URL is baked into the bundle at build time. Override per deploy:
#   docker build --build-arg VITE_API_URL=https://sentence.valc.cc/api ...
ARG VITE_API_URL=https://sentence.valc.cc/api
ENV VITE_API_URL=${VITE_API_URL}

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# --- runtime stage ---
FROM nginx:1.27-alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
