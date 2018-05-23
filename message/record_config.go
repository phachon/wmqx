package message

type RecordConfig struct {
	File *RecordFileConfig
}

func NewRecordConfigFile(config *RecordFileConfig) *RecordConfig {
	return &RecordConfig{
		File: config,
	}
}
