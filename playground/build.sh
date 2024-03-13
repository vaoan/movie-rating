docker-compose down -v --remove-orphans
docker-compose rm -f -s
docker-compose up --always-recreate-deps --remove-orphans --renew-anon-volumes --build