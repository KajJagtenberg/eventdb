FROM node:14-alpine AS build

WORKDIR /src

COPY package.json .
COPY yarn.lock .

RUN yarn

COPY . .

RUN yarn build

RUN yarn --production

FROM node:14-alpine

WORKDIR /app

COPY --from=build /src/package.json .
COPY --from=build /src/node_modules node_modules
COPY --from=build /src/.next .next

USER node

CMD ["yarn", "start"]