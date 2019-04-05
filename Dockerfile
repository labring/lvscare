FROM alpine:3.9.2
COPY lvscare /bin/lvscare
CMD ["lvscare"]
