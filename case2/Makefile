KIND_VERSION = 0.17.0
KUBERNETES_VERSION = 1.25.3

ROOT_DIR := $(CURDIR)
BIN_DIR := $(ROOT_DIR)/bin
KIND := $(BIN_DIR)/kind
KUBECTL := $(BIN_DIR)/kubectl

MANIFESTS_DIR := $(ROOT_DIR)/manifests

.PHONY: download-tools
download-tools:
	mkdir -p $(BIN_DIR)
	curl -sfL -o $(KIND) https://github.com/kubernetes-sigs/kind/releases/download/v$(KIND_VERSION)/kind-linux-amd64
	chmod a+x $(KIND)
	curl -sSL -o $(KUBECTL) https://dl.k8s.io/release/v$(KUBERNETES_VERSION)/bin/linux/amd64/kubectl
	chmod a+x $(KUBECTL)

.PHONY: create-cluster
create-cluster:
	$(KIND) create cluster --name=moco-issue-419 --image=kindest/node:v$(KUBERNETES_VERSION) --wait 1m

.PHONY: delete-cluster
delete-cluster:
	$(KIND) delete cluster --name=moco-issue-419

.PHONY: apply
apply:
	$(KUBECTL) apply -k $(MANIFESTS_DIR)
	$(KUBECTL) rollout status -w statefulsets/mysql

.PHONY: delete
delete:
	-$(KUBECTL) delete -k $(MANIFESTS_DIR)
	-$(KUBECTL) delete pvc data-mysql-0
	-$(KUBECTL) delete pvc data-mysql-1

.PHONY: init-source
init-source:
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "SET GLOBAL read_only = 0"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "CREATE USER 'repl'@'%' IDENTIFIED WITH mysql_native_password BY 'password'"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "GRANT REPLICATION SLAVE ON *.* TO 'repl'@'%'"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "INSTALL PLUGIN rpl_semi_sync_source SONAME 'semisync_source.so'"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "INSTALL PLUGIN rpl_semi_sync_replica SONAME 'semisync_replica.so'"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "SET GLOBAL rpl_semi_sync_source_enabled = 1"
	#$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "SET GLOBAL rpl_semi_sync_source_timeout = 86400000"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "SHOW VARIABLES LIKE 'rpl_semi_sync%'"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "FLUSH PRIVILEGES"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "SET GLOBAL super_read_only = 1"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "SELECT @@GLOBAL.gtid_executed, @@GLOBAL.read_only, @@GLOBAL.super_read_only"

.PHONY: init-replica
init-replica:
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "SET GLOBAL read_only = 0"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "INSTALL PLUGIN rpl_semi_sync_source SONAME 'semisync_source.so'"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "INSTALL PLUGIN rpl_semi_sync_replica SONAME 'semisync_replica.so'"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "SET GLOBAL rpl_semi_sync_replica_enabled = 1"
	#$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "SET GLOBAL rpl_semi_sync_source_timeout = 86400000"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "SHOW VARIABLES LIKE 'rpl_semi_sync%'"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "SET GLOBAL super_read_only = 1"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "SELECT @@GLOBAL.gtid_executed, @@GLOBAL.read_only, @@GLOBAL.super_read_only"

.PHONY: start-repl
start-repl:
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "STOP SLAVE IO_THREAD"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "CHANGE REPLICATION SOURCE TO SOURCE_HOST = 'mysql-0.mysql.default.svc', SOURCE_USER = 'repl', SOURCE_PASSWORD = 'password', SOURCE_AUTO_POSITION = 1"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "START SLAVE IO_THREAD"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "START REPLICA"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "SET GLOBAL read_only = 0"

.PHONY: show-source
show-source:
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "SELECT @@GLOBAL.server_id, @@GLOBAL.server_uuid, @@GLOBAL.gtid_executed, @@GLOBAL.read_only, @@GLOBAL.super_read_only\G"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "SELECT * FROM mysql.gtid_executed\G"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "SHOW MASTER STATUS\G"
	$(KUBECTL) exec -it mysql-0 -- mysql -u root -e "SHOW REPLICA STATUS\G"

.PHONY: show-replica
show-replica:
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "SELECT @@GLOBAL.server_id, @@GLOBAL.server_uuid, @@GLOBAL.gtid_executed, @@GLOBAL.read_only, @@GLOBAL.super_read_only\G"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "SELECT * FROM mysql.gtid_executed\G"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "SHOW MASTER STATUS\G"
	$(KUBECTL) exec -it mysql-1 -- mysql -u root -e "SHOW REPLICA STATUS\G"

.PHONY: exec-source
exec-source:
	$(KUBECTL) exec -it mysql-0 -- mysql -u root

.PHONY: exec-replica
exec-replica:
	$(KUBECTL) exec -it mysql-1 -- mysql -u root

.PHONY: exec
exec:
	$(KUBECTL) exec -it deploy/client -- bash

.PHONY: stop-source
stop-source:
	$(KUBECTL) exec mysql-0 -- touch /tmp/stop

.PHONY: start-source
start-source:
	$(KUBECTL) exec mysql-0 -- rm /tmp/stop

