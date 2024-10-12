PROJECT_DIR := $(CURDIR)/src
BUILD_DIR := $(CURDIR)/build
EXEC := Main
EXT := $(if $(findstring Windows_NT,$(OS)),.exe,)

# Commandes de construction et d'exécution
BUILD_CMD := go build -o $(BUILD_DIR)/$(EXEC)$(EXT)
RUN_CMD := $(BUILD_DIR)/$(EXEC)$(EXT)

# Préparation des modules Go
setup:
	cd $(PROJECT_DIR) && go mod tidy
	cd $(PROJECT_DIR) && go get github.com/gen2brain/raylib-go/raylib

# Construction du projet
build:
	cd $(PROJECT_DIR) && $(BUILD_CMD)

# Exécution du projet
run:
	$(RUN_CMD)

# Nettoyage des fichiers de build
clean:
	cd $(PROJECT_DIR) && rm -f $(BUILD_DIR)/$(EXEC)$(EXT) && go clean -modcache

# Cible pour préparer, construire et exécuter sans rien détruire
all: setup build run

.PHONY: setup build run clean all
