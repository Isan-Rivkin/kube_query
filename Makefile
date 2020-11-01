release:
	go build main.go
	mv main kq 
	cp kq /usr/local/bin/
test:
	kq get pods -n kube-system
clean: 
	rm /usr/local/bin/kq