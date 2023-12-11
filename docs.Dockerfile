FROM golang:1.21-alpine as base

# Load build arg
ARG APP_HOME

# Create development stage and set app location in container to /app
FROM base as development

WORKDIR $APP_HOME

# Add bash shell and gcc toolkit
RUN apk add --no-cache --upgrade bash build-base

# Add node/npm
RUN apk add --no-cache --update npm

# Add java (openapicli dep)
RUN apk add openjdk11-jre

# Install swaggo
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Install openapicli & redocly cli
RUN npm i -g @openapitools/openapi-generator-cli @redocly/cli@latest

# Copy go.mod and install go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Build and run binary with live reloading
CMD ["watch", "-n", "90", "make", "docs-gen"]
