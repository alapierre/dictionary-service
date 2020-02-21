CONNECTION_STRING=postgres://app:qwedsazxc@localhost:5432/app?sslmode=disable

modelgen:
	genna model -c $(CONNECTION_STRING) -o model/model.go -k
