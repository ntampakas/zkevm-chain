FROM golang:1.18 as preparego
COPY misc/prepare_go_env.sh /tmp
RUN chmod +x /tmp/prepare_go_env.sh
ENTRYPOINT ["/tmp/prepare_go_env.sh"]
