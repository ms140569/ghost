BIN_NAME = ghostd
TARGET_DIR=.
PARSER_DIR=parser
GLOBALS_DIR=globals
SERVER_DIR=server
ADMIN_DIR=server/admin
STORAGE_DIR=storage
CONSTANTS_DIR=constants
SCANNER_FILE=$(PARSER_DIR)/stomp.go
ADMIN_FILE=$(ADMIN_DIR)/adminshell.go

BASENAME=$(shell basename ${PWD})
NOW_STRING=$(shell date +%Y%m%d-%H%M)
BACKUP_FILE=$(BASENAME)-$(NOW_STRING).tar.gz

all: $(TARGET_DIR)/$(BIN_NAME)

TMP_FILES = $(SCANNER_FILE) $(ADMIN_FILE)

SRC = main.go $(PARSER_DIR)/command.go $(PARSER_DIR)/framebuilder.go $(PARSER_DIR)/frame.go $(PARSER_DIR)/token.go \
	$(GLOBALS_DIR)/config.go \
	$(GLOBALS_DIR)/configfile.go \
	$(SERVER_DIR)/server.go \
	$(SERVER_DIR)/session.go \
	$(SERVER_DIR)/heartbeat.go \
	$(SERVER_DIR)/frameprocessor.go \
	$(STORAGE_DIR)/storage.go \
	$(ADMIN_DIR)/admin.go \
	$(ADMIN_DIR)/adminapi.go \
	$(ADMIN_DIR)/shellcommand.go \
	$(STORAGE_DIR)/file.go \
	$(STORAGE_DIR)/memory.go \
	$(CONSTANTS_DIR)/constants.go \
	log/logger.go log/level/level.go

$(TARGET_DIR)/$(BIN_NAME): $(SCANNER_FILE) $(ADMIN_FILE) $(SRC)
	go build -o $(BIN_NAME)

.PHONY: clean
clean:
	@$(RM) -rf $(TARGET_DIR)/$(BIN_NAME) $(TMP_FILES)

.PHONY: fmt
fmt:
	go fmt

.PHONY: stat
stat: clean
	@find . -type f -name \*.go -o -name \*.py  |xargs wc -l

$(SCANNER_FILE): $(PARSER_DIR)/stomp.rl
	ragel -Z -T0 -o $(SCANNER_FILE) $(PARSER_DIR)/stomp.rl
	ragel -V -Z -p -o $(PARSER_DIR)/stomp.dot $(PARSER_DIR)/stomp.rl 

$(ADMIN_FILE): $(ADMIN_DIR)/adminshell.rl
	ragel -Z -T0 -o $(ADMIN_FILE) $(ADMIN_DIR)/adminshell.rl 
	ragel -V -Z -p -o $(ADMIN_DIR)/adminshell.dot $(ADMIN_DIR)/adminshell.rl

backup: clean
	(cd .. ; tar czvf $(BACKUP_FILE) $(BASENAME) ; cd -)

.PHONY: test
test: all
	(cd test;./run.sh)
	(cd test/generator;./runall.sh)
	(cd test/content-length;./run.sh)
	go test github.com/ms140569/ghost/log/level
	go test github.com/ms140569/ghost/parser
	go test github.com/ms140569/ghost/storage
	go test github.com/ms140569/ghost/server
	go test github.com/ms140569/ghost/server/admin

run: $(TARGET_DIR)/$(BIN_NAME)
	$(TARGET_DIR)/$(BIN_NAME)		
