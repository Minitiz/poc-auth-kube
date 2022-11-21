build-docker:
	docker build . -t api-test:local
	minikube image unload api-test:local
	minikube image load api-test:local

apply:
	kubectl config use-context minikube
	kubectl apply -f deployment.yaml

apply-api2:
	kubectl config use-context minikube
	kubectl apply -f deployment-only-get.yaml

apply-api3:
	kubectl config use-context minikube
	kubectl apply -f deployment-no-permission.yaml

clean:
	kubectl delete -f deployment.yaml
	kubectl delete -f deployment-only-get.yaml
	kubectl delete -f deployment-no-permission.yaml

pf:
	kubectl port-forward -n api svc/app  8080:8080

pf-api2:
	kubectl port-forward -n api2 svc/app  8080:8080

pf-api3:
	kubectl port-forward -n api3 svc/app  8080:8080

log:
	kubectl logs -n api -l app=api -f

log-api2:
	kubectl logs -n api2 -l app=api2 -f

log-api3:
	kubectl logs -n api3-l app=api3 -f