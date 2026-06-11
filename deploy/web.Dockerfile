FROM node:22-bookworm-slim AS deps

WORKDIR /app
COPY web/package.json web/package-lock.json ./
RUN npm ci

FROM node:22-bookworm-slim AS build

WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY web ./
RUN npm run build

FROM node:22-bookworm-slim AS runtime

WORKDIR /app
ENV NODE_ENV=production
COPY --from=build /app/.next/standalone ./
COPY --from=build /app/.next/static ./.next/static
COPY --from=build /app/public ./public
EXPOSE 3000
CMD ["node", "server.js"]
