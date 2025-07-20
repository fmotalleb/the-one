FROM gcr.io/distroless/base-debian12:nonroot AS distroless
COPY the-one /
ENTRYPOINT ["/the-one"]
