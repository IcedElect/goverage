.PHONY: profile cover

.profile:
	go test -coverprofile=coverage.out ./...

profile: .profile

cover-internal: .profile
	go tool cover -html=coverage.out

cover: .profile
	go run ../cmd/oh-my-cover-go --profile=coverage.out