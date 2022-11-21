build-docker:
	docker build -t storage-hub:local .
	minikube image unload storage-hub:local
	minikube image load storage-hub:local

apply:
	kubectl config use-context minikube
	kubectl apply -f deployment.yaml

clean:
	kubectl delete -f deployment.yaml

pf:
	kubectl port-forward -n storage-hub svc/app  8081:80

log:
	kubectl logs -n storage-hub -l app=storage-hub -f