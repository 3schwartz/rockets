messages:
	./lunar-backend-engineer-challenge/darwin_arm64/rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1

get:
	curl http://localhost:8088/messages/$(id)
