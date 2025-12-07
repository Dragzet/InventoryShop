# Build frontend with Vite
FROM node:18-alpine AS build
WORKDIR /app
COPY package.json package-lock.json* ./
COPY . .
RUN npm ci --silent || npm install --silent
RUN npm run build

FROM node:18-alpine AS runtime
WORKDIR /app
RUN npm install -g serve
COPY --from=build /app/dist /app/dist
EXPOSE 5173
CMD ["npx", "serve", "-s", "dist", "-l", "5173"]

