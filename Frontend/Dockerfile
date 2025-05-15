# Stage 1: Build the application
FROM node:18-alpine as builder

WORKDIR /app

# Copy package files first to leverage Docker cache
COPY package*.json ./

RUN npm ci --prefer-offline --no-audit

# Copy the rest of the files
COPY . .

# Build the application
RUN npm run build

# Stage 2: Serve the application
FROM node:18-alpine as production

WORKDIR /app

# Copy necessary files from builder
COPY --from=builder /app/.output /app/.output
COPY --from=builder /app/node_modules /app/node_modules
COPY --from=builder /app/package.json /app/package.json

# Fix tslib path issue by ensuring correct module structure
RUN if [ -f "/app/node_modules/tslib/tslib.js" ]; then \
      mkdir -p /app/node_modules/tslib/modules && \
      cp /app/node_modules/tslib/tslib.js /app/node_modules/tslib/modules/index.js; \
    fi

EXPOSE 3000
ENV NITRO_HOST=0.0.0.0
CMD ["node", "/app/.output/server/index.mjs"]