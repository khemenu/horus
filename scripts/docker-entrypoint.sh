hr init
hr create user admin \
	&& printf "admin\n" | hr --as admin set password \
	|| true

exec "$@"
