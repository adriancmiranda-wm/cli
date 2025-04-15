#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# ░░░░ Task Manager
#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# 📖 Referência: https://just.systems/man/en/
# 📖 Exemplos: https://github.com/casey/just/tree/master/examples
#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

set positional-arguments := true
set dotenv-load := true
set working-directory := "."

# 🏠 Alvo padrão
default: help

# 📋 Mostra este menu de ajuda
@help:
    just --list --unsorted \
      --list-heading $'🚀 WM - Available recipes…\n' \
      --list-prefix 'just '

#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# ⚙️ Setup
#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# 🔧 Instala ferramentas Go
setup-go:
    go install golang.org/x/tools/cmd/goimports@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install github.com/goreleaser/goreleaser/v2@latest
    go install golang.org/x/vuln/cmd/govulncheck@latest

# 🔧 Prepara o ambiente de desenvolvimento (instala ferramentas)
setup:
    #!/usr/bin/env bash
    # -*- coding: utf-8 -*-
    set -euo pipefail
    echo "📦 Instalando ferramentas auxiliares..."
    just setup-go
    echo "✅ Ambiente pronto!"
    echo
    BIN_PATH="${GOBIN:-${GOPATH:-${HOME}/go}/bin}";
    echo "🔍 Binários instalados em: ${BIN_PATH}";
    echo "ℹ️  Adicione ao PATH, se ainda não estiver:";
    printf "    export PATH=\$$PATH:%s\n" "${BIN_PATH}";

#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# 🔧 Configuração
#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# Nome do binário e caminho de build

PKG := "./..."
BIN := "bin/wm"
CLI := file_name(BIN)

#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# 🟢 Básico
#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# 🖨️ Exibe valores das variáveis de configuração
show-config:
    @echo "BIN: {{ BIN }}"
    @echo "CLI: {{ CLI }}"
    @echo "PKG: {{ PKG }}"

# 🧱 Compila o projeto localmente
build:
    go build -o {{ BIN }} .

# ▶️ Executa o projeto diretamente
run *ARGS:
    go run . {{ ARGS }}

# 🧹 Formata o código e ajusta imports
fmt:
    go fmt {{ PKG }}
    goimports -w .

# 🧩 Organiza dependências do módulo
deps:
    go mod tidy
    go mod verify

#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# 🟡 Intermediário
#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# 🧪 Executa testes com cobertura
test:
    go test -v -cover {{ PKG }}

# 🔍 Lint e checagens estáticas
lint:
    go vet {{ PKG }}
    golangci-lint run

# 🚀 Instala o binário globalmente (em $GOBIN)
install:
    go install .

# 💬 Injeta versão, commit e data e compila o binário
version-build VERSION="dev":
    #!/usr/bin/env bash
    # -*- coding: utf-8 -*-
    set -euo pipefail
    VERSION_VAL="{{ VERSION }}"
    COMMIT_VAL="$(git rev-parse --short HEAD)"
    DATE_VAL="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"

    # Lista de ldflags
    LDFLAGS=(
    	-s -w
    	-X main.Version=${VERSION_VAL}
    	-X main.Commit=${COMMIT_VAL}
    	-X main.Date=${DATE_VAL}
    )

    # Build com ldflags
    go build -ldflags="${LDFLAGS[*]}" -o {{ BIN }} .

#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# 🔵 Avançado
#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# 🌍 Compila para múltiplas plataformas (snapshot local, sem publicar)
release-local:
    goreleaser release --snapshot --clean --skip=publish --rm-dist

# 📦 Gera e publica uma versão real (requer GITHUB_TOKEN, usa a tag git)
release:
    goreleaser release --clean

# 🔒 Verifica vulnerabilidades conhecidas
audit:
    govulncheck {{ PKG }}

# 🧠 Benchmark de performance
bench:
    go test -bench=. -benchmem {{ PKG }}

#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# 📝 Autocompletes
#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# Gera autocompletes bash/zsh/fish
completions:
    mkdir -p completions
    just run completion bash > completions/{{ CLI }}.bash
    just run completion zsh > completions/_{{ CLI }}
    just run completion fish > completions/{{ CLI }}.fish

#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# 🧼 Utilitários
#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# 🗑️ Limpa binários e caches
clean:
    rm -rf bin/
    go clean -testcache -modcache

# 🧭 Mostra versão Go e módulo atual
info:
    go version
    go list -m all

#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# 🧱 Uso
#~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

init:
    just run template init --template golib --name MeuProjeto --author Adrian
