FROM node:lts-slim
WORKDIR /app
RUN apt-get update -y && apt-get install -y openssl

RUN npm i -g pnpm
COPY package.json pnpm-lock.yaml ./
RUN pnpm i

# generate prisma client
COPY prisma ./prisma
RUN pnpm prisma generate

# copy source code
COPY . .
RUN pnpm check
RUN pnpm build

CMD [ "node", "build/index.js" ]
