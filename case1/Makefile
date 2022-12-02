include Makefile.env

KUBERNETES_VERSION = 1.24.5
CERT_MANAGER_VERSION := 1.10.1
MOCO_VERSION := 0.14.0

ROOT_DIR := $(CURDIR)
BIN_DIR := $(ROOT_DIR)/bin
MANIFESTS_DIR := $(ROOT_DIR)/manifests
KUBECTL := $(BIN_DIR)/kubectl
KUBECTLMOCO := $(BIN_DIR)/kubectl-moco

export USE_GKE_GCLOUD_AUTH_PLUGIN=True

.PHONY: download-tools
download-tools:
	mkdir -p $(BIN_DIR)
	curl -sSL -o $(BIN_DIR)/kubectl https://dl.k8s.io/release/v$(KUBERNETES_VERSION)/bin/linux/amd64/kubectl
	chmod a+x $(BIN_DIR)/kubectl
	curl -sSL https://github.com/cybozu-go/moco/releases/download/v$(MOCO_VERSION)/kubectl-moco_v$(MOCO_VERSION)_linux_amd64.tar.gz \
	| tar xz -C $(BIN_DIR) kubectl-moco

.PHONY: create-cluster
create-cluster:
	gcloud container clusters create $(CLUSTER) --project $(PROJECT) --zone $(ZONE) \
	  --cluster-version $(KUBERNETES_VERSION) \
	  --machine-type n1-standard-32
	gcloud container clusters get-credentials $(CLUSTER) --project $(PROJECT) --zone $(ZONE)

.PHONY: delete-cluster
delete-cluster:
	gcloud container clusters delete $(CLUSTER) --project $(PROJECT) --zone $(ZONE)

.PHONY: kubeconfig
kubeconfig:
	gcloud container clusters get-credentials $(CLUSTER) --project $(PROJECT) --zone $(ZONE)

.PHONY: deploy-cert-manager
deploy-cert-manager:
	$(KUBECTL) apply -f https://github.com/jetstack/cert-manager/releases/download/v$(CERT_MANAGER_VERSION)/cert-manager.yaml
	$(KUBECTL) -n cert-manager wait --for=condition=available --timeout=180s --all deployments

.PHONY: deploy-moco
deploy-moco:
	curl -sSL -o $(MANIFESTS_DIR)/moco/upstream.yaml https://github.com/cybozu-go/moco/releases/download/v$(MOCO_VERSION)/moco.yaml
	$(KUBECTL) apply -k $(MANIFESTS_DIR)/moco
	$(KUBECTL) -n moco-system wait --for=condition=available --timeout=180s --all deployments

.PHONY: target
target:
	$(KUBECTL) apply -f $(MANIFESTS_DIR)/target.yaml
	$(KUBECTL) wait --for=condition=available --timeout=180s --all mysqlclusters
	$(KUBECTLMOCO) mysql -u moco-writable target -- -e "CREATE USER 'user1'@'%' IDENTIFIED BY 'pass'"
	$(KUBECTLMOCO) mysql -u moco-writable target -- -e "CREATE DATABASE db1"
	$(KUBECTLMOCO) mysql -u moco-writable target -- -e "GRANT ALL ON db1.* TO 'user1'@'%'"
	$(KUBECTLMOCO) mysql -u moco-writable target -- -e "CREATE TABLE db1.t1 (id int AUTO_INCREMENT, v1 int, INDEX(id));"
	$(KUBECTLMOCO) mysql -u moco-writable target -- -e "INSERT INTO db1.t1 (v1) VALUES (1); COMMIT;"

.PHONY: delete-target
delete-target: delete-client
	-$(KUBECTL) delete -f $(MANIFESTS_DIR)/target.yaml

.PHONY: client
client:
	$(KUBECTL) apply -f $(MANIFESTS_DIR)/client.yaml

.PHONY: delete-client
delete-client:
	-$(KUBECTL) delete -f $(MANIFESTS_DIR)/client.yaml

.PHONY: delete
delete: delete-client delete-target

.PHONY: logs
logs:
	$(KUBECTL) logs -n moco-system -l app.kubernetes.io/component=moco-controller -f

.PHONY: watch-pod
watch-pod:
	kubectl get pod -l app.kubernetes.io/created-by=moco -w

.PHONY: watch-cluster
watch-cluster:
	kubectl get mysqlcluster -w

.PHONY: show
show:
	$(KUBECTL) get mysqlclusters
	$(KUBECTL) top pod

.PHONY: count
count:
	$(KUBECTLMOCO) mysql target --index 0 -- -e "SELECT COUNT(*) FROM db1.t1\G"
	$(KUBECTLMOCO) mysql target --index 1 -- -e "SELECT COUNT(*) FROM db1.t1\G"
	$(KUBECTLMOCO) mysql target --index 2 -- -e "SELECT COUNT(*) FROM db1.t1\G"

.PHONY: gtid
gtid:
	$(KUBECTLMOCO) mysql target --index 0 -- -e "SELECT @@GLOBAL.gtid_executed, @@GLOBAL.super_read_only, @@GLOBAL.read_only\G"
	$(KUBECTLMOCO) mysql target --index 1 -- -e "SELECT @@GLOBAL.gtid_executed, @@GLOBAL.super_read_only, @@GLOBAL.read_only\G"
	$(KUBECTLMOCO) mysql target --index 2 -- -e "SELECT @@GLOBAL.gtid_executed, @@GLOBAL.super_read_only, @@GLOBAL.read_only\G"

.PHONY: processlist
processlist:
	$(KUBECTLMOCO) mysql target --index 0 -- -e "SHOW PROCESSLIST\G"
	$(KUBECTLMOCO) mysql target --index 1 -- -e "SHOW PROCESSLIST\G"
	$(KUBECTLMOCO) mysql target --index 2 -- -e "SHOW PROCESSLIST\G"
