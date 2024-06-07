FROM golang:alpine AS build

COPY . /build
WORKDIR /build
RUN go build .

FROM alpine
COPY --from=build /build/form-service /app/form-service
WORKDIR /app
CMD /app/form-service
