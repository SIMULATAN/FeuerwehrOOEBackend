FROM alpine:3.17

WORKDIR /app

ARG TARGETPLATFORM

COPY $TARGETPLATFORM/main .

RUN chmod +x main

ENTRYPOINT [ "/app/main" ]