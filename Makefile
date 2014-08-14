BIN_NAME = ghostd
TARGET_DIR=.
SCANNER_FILE=stomp.go
BASENAME=$(shell basename ${PWD})
NOW_STRING=$(shell date +%Y%m%d-%H%M)
BACKUP_FILE=$(BASENAME)-$(NOW_STRING).tar.gz

all: $(TARGET_DIR)/$(BIN_NAME)

TMP_FILES = $(SCANNER_FILE)
SRC = main.go command.go

$(TARGET_DIR)/$(BIN_NAME): $(SCANNER_FILE)
	go build -o $(BIN_NAME) $(SCANNER_FILE) $(SRC)

.PHONY: clean
clean:
	$(RM) -rf $(TARGET_DIR)/$(BIN_NAME) $(TMP_FILES)

$(SCANNER_FILE): stomp.rl
	ragel -Z -T0 -o $(SCANNER_FILE) stomp.rl 

backup: clean
	(cd .. ; tar czvf $(BACKUP_FILE) $(BASENAME) ; cd -)

.PHONY: test
test:
	(cd test;./run.sh)

run: $(TARGET_DIR)/$(BIN_NAME)
	$(TARGET_DIR)/$(BIN_NAME)		
