docker run --name redis -d bigzelda/redis
docker run --link redis:db --name bigzelda --publish 6060:8000 -d bigzelda/app