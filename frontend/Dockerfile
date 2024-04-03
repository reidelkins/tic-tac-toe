# Stage 1: Build the Angular application
FROM node:18 as build-step
WORKDIR /app
COPY package.json /app/
RUN npm install
COPY . /app
RUN npm run build

# Stage 2: Serve the application from Nginx
FROM nginx:alpine
COPY --from=build-step /app/dist/tic-tac-toe-front /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
