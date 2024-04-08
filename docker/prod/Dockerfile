FROM golang:1.22-bookworm AS builder
 
RUN apt-get update
RUN apt-get install -y ca-certificates curl gnupg
RUN mkdir -p /etc/apt/keyrings
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
RUN NODE_MAJOR=20
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_21.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list
RUN apt-get update
RUN apt-get install nodejs -y
RUN npm i -g nodemon webpack webpack-cli script-loader

 
COPY cmd /src/cmd
COPY components /src/components 
COPY core /src/core
COPY models /src/models
COPY pages /src/pages
COPY frontend /src/frontend
COPY static /src/static
COPY tailwind /src/tailwind
COPY tailwind.config.js /src/tailwind.config.js
COPY webpack.config.js /src/webpack.config.js

COPY go.mod /src/go.mod 
COPY go.sum /src/go.sum
  
WORKDIR /src 

RUN go build -o dist ./cmd
RUN webpack --mode=production
RUN ./tailwind/tailwindcss -i frontend/css/global.css -o static/.dist/styles.css --minify

FROM golang:1.22-bookworm

COPY --from=builder /src/dist /dist/dist
COPY --from=builder /src/static/ /dist/static/ 

WORKDIR /dist 
 
CMD [ "/dist/dist" ] 