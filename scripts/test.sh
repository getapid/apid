#/usr/local/bin sh

docker-compose -f tests/echo/docker-compose.yaml up -d

# positive tests first
failing=$($1 check -s "tests/**/*_pass.jsonnet" --json | jq ".failing")
if [ $failing -gt 0 ]; then
    echo "positive tests failed"
    exit 1
fi

# then negative tests
passing=$($1 check -s "tests/**/*_fail.jsonnet" --json | jq ".passing")
if [ $passing -gt 0 ]; then
    echo "negative tests failed"
    exit 1
fi

docker-compose -f tests/echo/docker-compose.yaml down