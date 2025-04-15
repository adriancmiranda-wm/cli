package templates

import (
	"embed"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed embedded/*
var embeddedFS embed.FS

type TemplateMeta struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Variables   []string `json:"variables,omitempty"`
	PostInit    []string `json:"post_init,omitempty"`
}

type TemplateFS struct {
	Meta  TemplateMeta
	Files map[string]string
}

// LoadTemplateFromDir lê um template de um diretório que contém template.json e arquivos *.tmpl
func LoadTemplateFromDir(dir string) (*TemplateFS, error) {
	metaFile := filepath.Join(dir, "template.json")

	data, err := os.ReadFile(metaFile)
	if err != nil {
		return nil, err
	}

	var meta TemplateMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	files := make(map[string]string)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".tmpl") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(dir, path)
			outPath := strings.TrimSuffix(rel, ".tmpl")
			files[outPath] = string(content)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &TemplateFS{Meta: meta, Files: files}, nil
}

// LoadEmbeddedTemplate carrega template embutido via embed.FS
func LoadEmbeddedTemplate(name string) (*TemplateFS, error) {
	base := "embedded/" + name
	metaPath := base + "/template.json"
	data, err := embeddedFS.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}
	var meta TemplateMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	files := make(map[string]string)
	entries, err := embeddedFS.ReadDir(base)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".tmpl") {
			content, _ := embeddedFS.ReadFile(filepath.Join(base, name))
			outPath := strings.TrimSuffix(name, ".tmpl")
			files[outPath] = string(content)
		}
	}

	return &TemplateFS{Meta: meta, Files: files}, nil
}

// RenderTemplate gera os arquivos do template no diretório de destino
func RenderTemplate(dest string, tpl *TemplateFS, vars map[string]string) error {
	if tpl == nil {
		return errors.New("template is nil")
	}
	for relPath, content := range tpl.Files {
		fullPath := filepath.Join(dest, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			return err
		}
		t, err := template.New(relPath).Parse(content)
		if err != nil {
			return err
		}
		f, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		if err := t.Execute(f, vars); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}
