# moco-issue-419

```bash
# Console (1)
cat << EOF > Makefile.env
PROJECT := <Your project>
ZONE := <Zone: e.g. asia-northeast1-c>
CLUSTER := moco-issue-419
EOF

make download-tools
make create-cluster deploy-cert-manager deploy-moco

# Console (2)
make watch-cluster

# Console (3)
make watch-pod

# Console (4)
make logs

# Console (1)
make target
# Wait until the cluster status becomes healthy.

make client
```
