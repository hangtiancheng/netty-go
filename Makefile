.PHONY: commit
commit: config
	git add -A && git commit -m 'feat: Introduce new features'

.PHONY: feat
feat: commit
	git push origin main

.PHONY: config
config:
	git config user.name hangtiancheng && git config user.email '161043261@qq.com'

AUTOBAHN_DIR := ws/autobahn

# Test the netty-go WebSocket server with Autobahn's fuzzing client.
.PHONY: autobahn-server
autobahn-server:
	mkdir -p $(AUTOBAHN_DIR)/report
	docker run -it --rm \
		-v "$(CURDIR)/$(AUTOBAHN_DIR)/config:/config" \
		-v "$(CURDIR)/$(AUTOBAHN_DIR)/report:/report" \
		crossbario/autobahn-testsuite \
		wstest -m fuzzingclient -s /config/fuzzingclient.json

# Test the netty-go WebSocket client with Autobahn's fuzzing server.
.PHONY: autobahn-client
autobahn-client:
	mkdir -p $(AUTOBAHN_DIR)/report
	docker run -it --rm \
		-v "$(CURDIR)/$(AUTOBAHN_DIR)/config:/config" \
		-v "$(CURDIR)/$(AUTOBAHN_DIR)/report:/report" \
		-p 9001:9001 \
		-p 9002:8080 \
		crossbario/autobahn-testsuite \
		wstest -m fuzzingserver -s /config/fuzzingserver.json
