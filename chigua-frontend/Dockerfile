# 构建阶段
FROM node:22 AS builder

# 设置工作目录
WORKDIR /app

# 复制package.json和pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./
# 安装pnpm
RUN npm config set registry https://registry.npmmirror.com
RUN npm install -g pnpm
# 安装依赖
RUN pnpm install

# 复制项目文件
COPY . .

# 构建项目
RUN pnpm run build

# 运行阶段
FROM nginx:1.29

# 复制构建产物到nginx的html目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置文件
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]
