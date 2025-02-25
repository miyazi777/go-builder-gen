package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"test1/constant"
)

type BuilderFile struct {
	path              string
	lines             []string
	beforeCommentLies []string
	afterCommentLines []string
}

func NewBuilderFile(path string) *BuilderFile {
	return &BuilderFile{path: path}
}

func (t *BuilderFile) IsExistFile() bool {
	outputFile := t.GetTargetFilePath()
	_, err := os.Stat(outputFile)
	return !os.IsNotExist(err)
}

func (t *BuilderFile) GetTargetFilePath() string {
	return filepath.Join(filepath.Dir(t.path), t.GetBuilderFileName(t.path))
}

func (t *BuilderFile) GetTargetDir() string {
	return filepath.Dir(t.path)
}

func (t *BuilderFile) GetBuilderFileName(targetFilePath string) string {
	baseNameWithExt := filepath.Base(targetFilePath)
	baseName := strings.TrimSuffix(baseNameWithExt, filepath.Ext(baseNameWithExt))
	return fmt.Sprintf("%s_builder.go", baseName)
}

func (t *BuilderFile) IsInvalidComment() error {
	if len(t.lines) == 0 {
		t.readFile()
	}

	count := 0
	for _, line := range t.lines {
		if strings.TrimSpace(line) == constant.COMMENT {
			count++
		}
	}

	if count != 2 {
		return errors.New(fmt.Sprintf("正しいコメント(%s)が2つ存在しない為、コード生成に失敗しました。\n", constant.COMMENT))
	}

	return nil
}

func (t *BuilderFile) GetBeforeCommentLines() string {
	if len(t.beforeCommentLies) == 0 {
		t.splitLines()
	}
	return strings.Join(t.beforeCommentLies, "\n")
}

func (t *BuilderFile) GetAfterCommentLines() string {
	if len(t.afterCommentLines) == 0 {
		t.splitLines()
	}
	return strings.Join(t.afterCommentLines, "\n")
}

func (t *BuilderFile) splitLines() {
	if len(t.lines) == 0 {
		t.readFile()
	}

	count := 0
	beforeLines := []string{}
	afterLines := []string{}
	for _, line := range t.lines {
		if strings.TrimSpace(line) == constant.COMMENT {
			count++
			continue
		}
		switch count {
		case 0:
			beforeLines = append(beforeLines, line)
		case 2:
			afterLines = append(afterLines, line)
		}
	}

	t.beforeCommentLies = beforeLines
	t.afterCommentLines = afterLines
}

func (t *BuilderFile) readFile() error {
	file, err := os.Open(t.GetTargetFilePath())
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	t.lines = []string{}
	for scanner.Scan() {
		t.lines = append(t.lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
