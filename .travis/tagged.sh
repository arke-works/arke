CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o arke-forum .;
docker build -t arkeworks/arke:$TRAVIS_TAG .;
docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD";
docker push arkeworks/arke:$TRAVIS_TAG;
bash <(curl -s https://raw.githubusercontent.com/goreleaser/get/master/latest);