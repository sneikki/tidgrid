FROM node:18.14.2-bullseye-slim

WORKDIR /usr/src/tidgrid
COPY package.json package-lock.json ./

RUN npm ci --omit=dev
COPY . .
RUN npm run build

CMD ["ls", "build"]
