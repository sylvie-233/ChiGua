FROM node:22 AS builder

WORKDIR /app

COPY package.json pnpm-lock.yaml ./
RUN npm config set registry https://registry.npmmirror.com
RUN npm install -g pnpm
RUN pnpm config set dangerouslyAllowAllBuilds true
RUN pnpm install --no-frozen-lockfile --verbose --unsafe-perm

COPY . .

RUN pnpm run build

FROM nginx:1.29

COPY --from=builder /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
