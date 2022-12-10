COVER_FILE = ".coverage.out"

cover:
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile=${COVER_FILE}
