# make startProm
.PHONY: startProm
startProm:
	docker run \
	--rm \
	-p 9090:9090 \
	--name=prometheus \
	-v $(shell pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
	prom/prometheus

# make startGrafana
# for first timers, the username & password is both `admin`
.PHONY: startGrafana
startGrafana:
	docker volume create grafana-storage
	docker volume inspect grafana-storage
	docker run -p 3000:3000 --name=grafana grafana/grafana-oss || docker start grafana

# make docker compose build
.PHONY: dockerComposeBuild
dockerComposeBuild:
	docker-compose build

# make docker compose up
.PHONY: dockerComposeUp
dockerComposeUp:
	docker-compose up

# make docker compose down
.PHONY: dockerComposeDown
dockerComposeDown:
	docker-compose down

# make docker compose restart
.PHONY: dockerComposeRestart
dockerComposeRestart:
	docker-compose restart

# make docker compose logs
.PHONY: dockerComposeLogs
dockerComposeLogs:
	docker-compose logs -f

# make docker compose force rebuild
.PHONY: dockerComposeForceRebuild
dockerComposeForceRebuild:
	docker-compose down
	docker-compose build
	docker-compose up
	