# Run

```bash
go build ./

./docker-implementation-practice run -ti /bin/sh
```

```bash with cgroup
./dip run -ti -m 100m stress --vm-bytes 200m --vm-keep -m 1
```