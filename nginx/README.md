docker の nginx で
/var/log/nginx/access.log
は stdout へのシンボリックリンクになっているので、
access_log を/var/log/nginx/access2.log へ出すようにする

alp はcheck.shで実行する
