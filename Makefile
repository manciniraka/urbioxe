# mock
mock-install:
	go install github.com/golang/mock/mockgen@latest
mock-report:
	mockgen -source internal/repository/report_repository.go -destination mocks/report_mock_repo.go
mock-staff:
	mockgen -source internal/repository/staff_repository.go -destination mocks/staff_mock_repo.go
mock-cloudinary:
	mockgen -source external/cloudinary/claudinary.go -destination mocks/cloudinary_mock_service.go

#test
test-reportSvc:
	go test -v ./internal/service/... -coverprofile=report-coverage.out -cover -failfast
test-reportSvc-coverage:
	go test -v $$(go list ./internal/service/... | grep -v '/mock') -coverprofile=report-coverage.out -cover -failfast && \
	go tool cover -html=report-coverage.out -o report-cover.html && \
	open report-cover.html
