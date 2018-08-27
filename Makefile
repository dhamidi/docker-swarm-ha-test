
logs:
	sudo docker service logs -f ha_http-service

watch:
	while true; do sleep 1; date; make --quiet health; done

health:
	@gateway_bridge=$$(sudo docker network inspect docker_gwbridge --format '{{(index .IPAM.Config 0).Gateway}}'); \
	curl "$$gateway_bridge":12345/health
naive:
	sudo docker stack deploy -c naive.yml ha

healthcheck:
	sudo docker stack deploy -c healthcheck.yml ha

build:
	make -C slow-starting-service

destroy:
	sudo docker stack rm ha
