# isucon2020-himetex

## 起動
sudo systemctl start iscon

## 自動起動
* edit tools/iscon.service

```shell script
mv tools/iscon.service /etc/systemd/system/iscon.service
```

sudo systemctl enable iscon