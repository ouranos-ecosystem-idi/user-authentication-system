#! /bin/bash

while read -r f; do
	PGPASSWORD="passw0rd" psql --set ON_ERROR_STOP=1 -h 127.0.0.1 -p 5432 -U dhuser -d dhlocal -f $f 1>/dev/null
	if [[ 0 -ne ${?} ]]; then
		exit 1
	fi
done < <(find setup/seeders -name "*.sql" -mindepth 1 -maxdepth 1 | sort)
