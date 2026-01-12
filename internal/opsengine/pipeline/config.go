package pipeline

// PipelineConfig 对应 .devnexus.yaml 的顶层结构
type PipelineConfig struct {
	Name   string  `yaml:"name"`   // 流水线名字
	Stages []Stage `yaml:"stages"` // 包含哪些阶段
}

type Stage struct {
	Name   string   `yaml:"name"`   // 阶段名称
	Image  string   `yaml:"image"`  // TODO：指定用哪个Docker镜像跑
	Script []string `yaml:"script"` // 要执行的Shell命令列表
}
