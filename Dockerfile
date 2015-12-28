FROM scratch
#
#   Docker container build for golang
#
ADD https://raw.githubusercontent.com/bagder/ca-bundle/master/ca-bundle.crt /etc/ssl/certs/ca-bundle.crt
ADD dozy /usr/local/bin/dozy
ENTRYPOINT ["/usr/local/bin/dozy"]
