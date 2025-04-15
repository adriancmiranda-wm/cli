# WM CLI (hybrid templates skeleton)

> CLI para gerenciamento geral de projetos

Este esqueleto implementa uma CLI em Go com suporte híbrido a templates:

- **embedded** via `embed.FS` (pequenos templates de arquivo/config)
- **defaults** editáveis dentro do repo (`internal/templates/defaults/*`) para projetos vazios
- **local** em `~/.wm/templates/<name>` — para templates de bibliotecas maiores
- suporte a `--from` para apontar um caminho local (ou futuro clone remoto)

## Como usar

1. Build

```bash
just build
# ou
go build -o bin/wm .
```

2. Listar templates disponíveis

```bash
./bin/wm template list
```

3. Gerar projeto usando default (editable)

```bash
./bin/wm template init -t golib -n mylib -a "Adrian"
```

4. Gerar projeto usando npm template (bun)

```bash
./bin/wm template init -t npm-package -n mypkg -a "Adrian"
```

5. Gerar arquivo simples a partir do template embutido

```bash
./bin/wm template init -t go-file -n MyFile -a "Adrian"
# go-file é embutido e gerará um único arquivo Example.go dentro de ./MyFile
```

6. Adicionar um template local (exemplo)

```bash
mkdir -p ~/.wm/templates/myproject
# coloque template.json e arquivos *.tmpl dentro
./bin/wm template init -t myproject -n Whatever -a "You"
```

## Estrutura do projeto

- cmd/: comandos CLI
- internal/templates/: lógica de templates + defaults editáveis e embedded
- internal/log/: logger e recuperação de panic
- justfile e .goreleaser.yaml prontos para uso

---

[//]: https://github.com/vercel/turborepo/tree/ff4c410966b784243f9a782b2f1106a0fa0605c8/packages/turbo-gen
