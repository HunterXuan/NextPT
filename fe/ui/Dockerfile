# Multi-stage build for minimal image size

# Stage 1: Build
FROM node:24-alpine AS builder

WORKDIR /app

# Install build dependencies for native modules
RUN apk add --no-cache python3 make g++

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the application with increased memory limit
ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN npm run build

# Stage 2: Production
FROM node:24-alpine AS production

WORKDIR /app

# Copy built application from builder
COPY --from=builder /app/.output /app/.output

# Expose port
EXPOSE 3000

# Set environment to production
ENV NODE_ENV=production

# Run the application
CMD ["node", ".output/server/index.mjs"]
