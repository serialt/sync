## Sync 下载同步Github release包

### 1、配置文件
配置文件默认读取`./config.yaml`,由于需要长时间访问github api，需要配置github access-token防止请求受限。


### 配置文件格式
```shell
[root@tc tmp]# cat ~/.sync.yaml 
# 频繁访问可能会导致github限流,可以配置多个不同账号的token和随机sleep
githubToken: 
  - "ghp_3PpD0000000000000000000000000000000"
  - "ghp_3PpD0000000000000000000000022222222"
  - "ghp_3PpD0000000000000000000000055555555"
randomSleep: 5
lastNum: 3
mirrorRoot: ./tmp
githubRelease:
  - "prometheus/prometheus"
  - "prometheus/mysqld_exporter"
  - "prometheus/alertmanager"
  - "prometheus/prometheus"
  - "prometheus/mysqld_exporter"
  - "prometheus/alertmanager"
  - "prometheus/haproxy_exporter"
  - "prometheus/node_exporter"
  - "prometheus/blackbox_exporter"
  - "prometheus/jmx_exporter"
  - "prometheus/consul_exporter"
  - "prometheus/snmp_exporter"
  - "prometheus/memcached_exporter"
  - "prometheus/pushgateway"
  - "prometheus/statsd_exporter"
  - "prometheus/influxdb_exporter"
  - "prometheus/collectd_exporter"
  - "ClickHouse/clickhouse_exporter"
  - "danielqsj/kafka_exporter"
  - "oliver006/redis_exporter"
  - "prometheus-community/elasticsearch_exporter"
  - "prometheus-community/windows_exporter"
  - "prometheus-community/postgres_exporter"
  - "prometheus-community/elasticsearch_exporter"
  - "prometheus-community/pgbouncer_exporter"
  - "prometheus-community/bind_exporter"
  - "prometheus-community/smartctl_exporter"
  - "percona/mongodb_exporter"
  - "iamseth/oracledb_exporter"
  - "ncabatoff/process-exporter"
  - "nginxinc/nginx-prometheus-exporter"
  - "cloudflare/ebpf_exporter"
  - "martin-helmich/prometheus-nginxlog-exporter"
  - "hnlq715/nginx-vts-exporter"
  - "vvanholl/elasticsearch-prometheus-exporter"
  - "free/sql_exporter"
  - "hipages/php-fpm_exporter"
  - "digitalocean/ceph_exporter"
  - "pryorda/vmware_exporter"
  - "Lusitaniae/apache_exporter"
  - "joe-elliott/cert-exporter"   
  - "v2fly/v2ray-core"
  - "fatedier/frp"
  - "v2rayA/v2rayA"
  - "XTLS/Xray-core"
  - "Eugeny/tabby"
# go: https://go.dev/dl/go1.17.7.darwin-amd64.tar.gz
# node: https://registry.npmmirror.com/-/binary/node/v16.14.0/node-v16.14.0-linux-x64.tar.gz 
# nginx: https://repo.huaweicloud.com/nginx/
# git-win: https://repo.huaweicloud.com/git-for-windows/
# git-mac: https://repo.huaweicloud.com/git-for-macos/
# helm: https://repo.huaweicloud.com/helm/
# grafana: https://repo.huaweicloud.com/grafana/
# 含有一下关键字的不下载
excludeTxt:
  - "darwin-amd64"
  - "linux-arm64"
  - "windows-arm64"
  - "390"
  - "386"
  - "mips"
  - "bsd"
  - "dragonfly"
  - "32"
  - "ppc"
  - "riscv"
  - "s3960"
  - "illumos"
  - "v5"
  - "v6"
  - "v7"
  - "Source"
  - "yml"
  - "yaml"
  - "pacman"
  - "portbale"
  - "blockmap"
```
