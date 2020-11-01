BIN_DIR='/usr/local/bin'

release:
	go build main.go
	mv main kq
	echo "Installing to ${BIN_DIR}" 
	mv kq ${BIN_DIR}
test:
	kq get pods -n kube-system
clean: 
	rm ${BIN_DIR}/kq