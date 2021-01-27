FROM node:alpine AS build

WORKDIR /app

COPY package.json .
COPY yarn.lock .

RUN yarn

COPY . .

RUN yarn build
RUN yarn export

FROM node:alpine

WORKDIR /app

COPY package.json .
COPY yarn.lock .

RUN yarn --production

COPY --from=build /app/.next ./.next

CMD ["yarn","start"]