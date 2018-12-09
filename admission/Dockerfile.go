FROM scratch
EXPOSE 443
ENTRYPOINT ["/webhook"]
COPY admission /webhook

# GOOS=linux CGO_ENABLED=0 go build .
# docker build -t lizrice/admission:0.2 .