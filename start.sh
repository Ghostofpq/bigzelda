docker run --name redis -d bigzelda/redis
docker run --link redis:db --name bigzelda --publish 80:8000 -d bigzelda/app