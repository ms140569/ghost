BIN_NAME = ghostd
TARGET_DIR=.
PARSER_DIR=parser
SCANNER_FILE=$(PARSER_DIR)/stomp.go
BASENAME=$(shell basename ${PWD})
NOW_STRING=$(shell date +%Y%m%d-%H%M)
BACKUP_FILE=$(BASENAME)-$(NOW_STRING).tar.gz

all: $(TARGET_DIR)/$(BIN_NAME)

TMP_FILES = $(SCANNER_FILE)
SRC = main.go $(PARSER_DIR)/command.go $(PARSER_DIR)/framebuilder.go $(PARSER_DIR)/frame.go $(PARSER_DIR)/token.go constants/constants.go server/server.go

$(TARGET_DIR)/$(BIN_NAME): $(SCANNER_FILE) $(SRC)
	go build -o $(BIN_NAME)

.PHONY: clean
clean:
	$(RM) -rf $(TARGET_DIR)/$(BIN_NAME) $(TMP_FILES)

.PHONY: fmt
fmt:
	go fmt

$(SCANNER_FILE): $(PARSER_DIR)/stomp.rl
	ragel -Z -T0 -o $(SCANNER_FILE) $(PARSER_DIR)/stomp.rl 

backup: clean
	(cd .. ; tar czvf $(BACKUP_FILE) $(BASENAME) ; cd -)

.PHONY: test
test:
	(cd test;./run.sh)

run: $(TARGET_DIR)/$(BIN_NAME)
	$(TARGET_DIR)/$(BIN_NAME)		
