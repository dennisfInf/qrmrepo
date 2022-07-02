#stage 1
FROM node:latest as node
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build
#stage 2
FROM nginx:alpine
COPY --from=node /app/dist/elon-mask /usr/share/nginx/html
COPY nginx.conf.template /etc/nginx/nginx.conf.template
CMD ["/bin/sh",  "-c",  "envsubst '${API_HOST}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf && envsubst < /usr/share/nginx/html/assets/env.sample.js > /usr/share/nginx/html/assets/env.js && exec nginx -g 'daemon off;'"]
