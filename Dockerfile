FROM alpine:latest

ADD cdp-volumes-customizer /cdp-volumes-customizer
ENTRYPOINT ["./cdp-volumes-customizer"]