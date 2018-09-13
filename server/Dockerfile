FROM scratch
ADD ../ca-certificates.crt /etc/ssl/certs/
ADD bin/deploy-server /
CMD ["/deploy-server"]
EXPOSE 8000