#stage 1
FROM node:latest as node
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build
#stage 2
FROM nginx:alpine
COPY --from=node /app/dist/elon-mask /usr/share/nginx/html
#COPY nginx.conf /etc/nginx/nginx.conf
CMD ["/bin/sh",  "-c",  "envsubst '${API_HOST}' < nginx.conf > /etc/nginx/nginx.conf && exec nginx -g 'daemon off;'"]
