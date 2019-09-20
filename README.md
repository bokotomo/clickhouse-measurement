
・clickhouseアクセス
```
cd docker && docker-compose up -d && docker-compose exec client bash

/usr/bin/clickhouse-client --host server --multiline
```


・容量チェック
```
cd /var/lib/clickhouse/data/

du -hs default
```
