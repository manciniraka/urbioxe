# mock
mock-install:
	go install github.com/golang/mock/mockgen@latest
mock-report:
	mockgen -source internal/repository/report_repository.go -destination mocks/report_mock_repo.go
mock-staff:
	mockgen -source internal/repository/staff_repository.go -destination mocks/staff_mock_repo.go
mock-cloudinary:
	mockgen -source external/cloudinary/claudinary.go -destination mocks/cloudinary_mock_service.go
mock-category:
	mockgen -source internal/repository/category_repository.go -destination mocks/category_mock_repo.go
mock-department:
	mockgen -source internal/repository/department_repository.go -destination mocks/department_mock_repo.go
mock-district:
	mockgen -source internal/repository/district_repository.go -destination mocks/district_mock_repo.go
mock-emergency:
	mockgen -source internal/repository/emergency_repository.go -destination mocks/emergency_mock_repo.go
mock-news:
	mockgen -source internal/repository/news_repository.go -destination mocks/news_mock_repo.go
mock-user:
	mockgen -source internal/repository/user_repository.go -destination mocks/user_mock_repo.go
mock-water:
	mockgen -source internal/repository/water_repository.go -destination mocks/water_mock_repo.go
mock-meter:
	mockgen -source internal/repository/meter_reading_repository.go -destination mocks/meter_reading_mock_repo.go

#test
test-service:
	go test -v ./internal/service/... -coverprofile=report-coverage.out -cover -failfast
test-service-coverage:
	go test -v $$(go list ./internal/service/... | grep -v '/mock') -coverprofile=report-coverage.out -cover -failfast && \
	go tool cover -html=report-coverage.out -o report-cover.html && \
	open report-cover.html
