messages:
	./lunar-backend-engineer-challenge/darwin_arm64/rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1

rocket:
	curl http://localhost:8088/rocket/$(id)

filterExploded ?= false
limit ?= -1

rockets:
	@curl "http://localhost:8088/rockets?filterExploded=$(filterExploded)&limit=$(limit)"

test:
	go test $$(go list ./...) -coverprofile cover.out
