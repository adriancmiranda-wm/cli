package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/adriancmiranda-wm/cli/internal/templates"
	"github.com/spf13/cobra"
)

var (
	projectName string
	tplName     string
	author      string
	fromRepo    string
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Gerencia templates de projetos",
}

var templateInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Gera um projeto a partir de template",
	RunE: func(cmd *cobra.Command, args []string) error {
		if projectName == "" {
			return fmt.Errorf("é necessário informar o nome do projeto com --name")
		}

		var tpl *templates.TemplateFS
		var err error

		// 1) if --from provided (local path) try it first
		if fromRepo != "" {
			if _, statErr := os.Stat(fromRepo); statErr == nil {
				tpl, err = templates.LoadTemplateFromDir(fromRepo)
				if err != nil {
					return err
				}
			} else {
				return errors.New("path provided in --from not found; remote clone is not implemented in this skeleton")
			}
		} else {
			// 2) look in ~/.wm/templates/<tplName>
			home, _ := os.UserHomeDir()
			localDir := filepath.Join(home, ".wm", "templates", tplName)
			if _, err := os.Stat(localDir); err == nil {
				tpl, err = templates.LoadTemplateFromDir(localDir)
				if err != nil {
					return err
				}
			} else {
				// 3) try internal defaults editable dirs
				defaultDir := filepath.Join("internal", "templates", "defaults", tplName)
				if _, err := os.Stat(defaultDir); err == nil {
					tpl, err = templates.LoadTemplateFromDir(defaultDir)
					if err != nil {
						return err
					}
				} else {
					// 4) fallback to embedded templates
					tpl, err = templates.LoadEmbeddedTemplate(tplName)
					if err != nil {
						return fmt.Errorf("template '%s' não encontrado em locais locais, defaults ou embedded", tplName)
					}
				}
			}
		}

		vars := map[string]string{
			"ProjectName": projectName,
			"Author":      author,
		}

		dest := filepath.Join(".", projectName)
		if err := templates.RenderTemplate(dest, tpl, vars); err != nil {
			return err
		}

		slog.Info("Projeto criado", "path", dest, "template", tplName)
		fmt.Println("✅ Projeto gerado em:", dest)
		return nil
	},
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista templates disponíveis (local, defaults, embedded)",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, _ := os.UserHomeDir()
		localRoot := filepath.Join(home, ".wm", "templates")
		fmt.Println("Available templates:")
		// list local
		if entries, err := os.ReadDir(localRoot); err == nil {
			for _, e := range entries {
				if e.IsDir() {
					fmt.Println("local:", e.Name())
				}
			}
		}
		// defaults
		defaultRoot := filepath.Join("internal", "templates", "defaults")
		if entries, err := os.ReadDir(defaultRoot); err == nil {
			for _, e := range entries {
				if e.IsDir() {
					fmt.Println("default:", e.Name())
				}
			}
		}
		// embedded (hardcoded list)
		fmt.Println("embedded: go-file, config")
		return nil
	},
}

func init() {
	templateInitCmd.Flags().StringVarP(&projectName, "name", "n", "", "Nome do projeto")
	templateInitCmd.Flags().StringVarP(&tplName, "template", "t", "golib", "Nome do template")
	templateInitCmd.Flags().StringVarP(&author, "author", "a", "Autor", "Nome do autor")
	templateInitCmd.Flags().StringVar(&fromRepo, "from", "", "Caminho local do template ou URL (clone remoto não implementado)")

	templateCmd.AddCommand(templateInitCmd)
	templateCmd.AddCommand(templateListCmd)
	rootCmd.AddCommand(templateCmd)
}
