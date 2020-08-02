DST_PATH = ./build/bin
ENTRY_POINT_PATH= ./cmd/mobilebanking
SCRIPTS_PATH = ./scripts
MIGRATION_SCRIPTS_PATH= ${SCRIPTS_PATH}/postgres/migration
INTERNAL_APP_PATH = ./internal/app
TEST_MOCK_PATH =  ./test/mock
.PHONY= mock
mock:
	mockgen -source=${INTERNAL_APP_PATH}/domain/user.go -destination=${TEST_MOCK_PATH}/mock_domain/user_mock.go
	mockgen -source=${INTERNAL_APP_PATH}/domain/savings.go -destination=${TEST_MOCK_PATH}/mock_domain/savings_mock.go
	mockgen -source=${INTERNAL_APP_PATH}/domain/transaction.go -destination=${TEST_MOCK_PATH}/mock_domain/transaction_mock.go
.PHONY= setuppostgres
setuppostgres:
	sh ${SCRIPTS_PATH}/postgres/setup-postgres.sh
.PHONY= migrationinit
initmigration:
	go run ${MIGRATION_SCRIPTS_PATH}/*.go init
.PHONY= migratetable
migratetable:
	go run ${MIGRATION_SCRIPTS_PATH}/*.go up 1
.PHONY= seeddata
seeddata:
	go run ${MIGRATION_SCRIPTS_PATH}/*.go up 2
.PHONY = resetdata
resetdata:
	go run ${MIGRATION_SCRIPTS_PATH}/*.go reset
.PHONY = setupdata
setupsdata: initmigration resetdata migratetable seeddata
	echo "good to go..."
.PHONY= e2e
e2e:initmigration resetdata migratetable seeddata
	go test -timeout 1000s github.com/tunaiku/mobilebanking/test/e2e -count=1 -v 
.PHONY= buildapp
buildapp: 
	rm -rf ${DST_PATH}
	go build -o ${DST_PATH}/mobile-banking-service ${ENTRY_POINT_PATH}/*.go
.PHONY= run
run: 
	go run ${ENTRY_POINT_PATH}/*.go
