// levels.yaml embed + 解析。
// 通过 go:embed 把 yaml 编进 binary,运行时不依赖外部文件。
package main

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed assets/levels.yaml
var levelsYAML []byte

func LoadLevels() ([]Level, error) {
	var data struct {
		Levels []Level `yaml:"levels"`
	}
	if err := yaml.Unmarshal(levelsYAML, &data); err != nil {
		return nil, fmt.Errorf("parse levels.yaml: %w", err)
	}
	if err := FinalizeLevels(data.Levels); err != nil {
		return nil, err
	}
	return data.Levels, nil
}
