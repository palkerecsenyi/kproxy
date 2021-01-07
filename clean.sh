docker-compose down

docker run --it \
  --env KPROXY_DB_PATH=/opt/db \
  --env KPROXY_PATH=/opt/cache \
  --env KPROXY_MAX_SPACE=536870912000 \
  -v /opt/db:/opt/db \
  -v /opt/cache:/opt/cache \
  kproxy_proxy go run main.go -clean

docker-compose up