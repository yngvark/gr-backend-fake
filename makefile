LOCAL_IMAGE = k8s-app-template

.PHONY: help
help: ## Print this menu
	@grep -E '^[a-zA-Z_0-9-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## -
	docker build . -t $(LOCAL_IMAGE)

push: build ##
	$(eval TAGGED_IMAGE := "docker.pkg.github.com/yngvark/k8s-app-template/k8s-app-template:latest")
	docker tag $(LOCAL_IMAGE) $(TAGGED_IMAGE)
	docker push $(TAGGED_IMAGE)

run: ## -
	docker run --name k8s-app-template --rm -p 8080:8080 k8s-app-template

init: ##
	@echo How to get docker password from github:
	@echo https://docs.github.com/en/free-pro-team@latest/packages/using-github-packages-with-your-projects-ecosystem/configuring-docker-for-use-with-github-packages
	@echo
	@echo Fill in "--docker-password" argument, and run:
	@echo kubectl create secret docker-registry regsecret --docker-server=docker.pkg.github.com --docker-username=yngvark --docker-password=MY_PASSWORD

deploy: ## -
	kubectl apply --namespace app-template -f k8s

undeploy: ## -
	kubectl delete --namespace app-template -f k8s
