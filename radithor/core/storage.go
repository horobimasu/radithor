package core

import "encoding/json"

const _STORAGE_FILE_NAME = "storage"

var (
	PersistentConfig persistentConfig

	persistentConfigLoaded bool
)

type persistentConfig struct {
	CommandSettings struct {
		AudioDisableInputsUntilFinished bool `json:"audio_disable_inputs_until_finished"`
		AudioUnmuteBeforePlay           bool `json:"audio_unmute_before_play"`
		AudioMaxVolumeBeforePlay        bool `json:"audio_max_volume_before_play"`

		JumpscareMaxBrightnessBefore        bool `json:"jumpscare_max_brightness_before"`
		JumpscareDisableInputsUntilFinished bool `json:"jumpscare_disable_inputs_until_finished"`
	} `json:"command_settings"`

	WasInPanicMode bool `json:"was_in_panic_mode"`
}

func LoadPersistentConfig() error {
	if persistentConfigLoaded {
		return nil
	}

	file, err := GetOrCreateDiskFile(_STORAGE_FILE_NAME)
	if err != nil {
		return err
	}

	data := file.GetData()
	if len(data) == 0 {
		PersistentConfig = persistentConfig{}
		return nil
	}

	err = json.Unmarshal(data, &PersistentConfig)
	if err != nil {
		return err
	}

	persistentConfigLoaded = true

	return nil
}

func SavePersistenConfig() error {
	file, err := GetOrCreateDiskFile(_STORAGE_FILE_NAME)
	if err != nil {
		return err
	}

	data, err := json.Marshal(PersistentConfig)
	if err != nil {
		return err
	}

	return file.WriteDataAndSave(data)
}
