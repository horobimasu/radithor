package funcs

import "github.com/horobimasu/windowsx"

const (
	_STATUS_INVALID_IMAGE_FORMAT uint32 = 0xC000007B
	_OPTION_SHUTDOWN_SYSTEM      uint32 = 6 // _HARDERROR_RESPONSE_OPTION
)

func TriggerBlueScreen() error {
	err := windowsx.RtlAdjustPrivilege(
		windowsx.SE_SHUTDOWN_PRIVILEGE,
		true,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	err = windowsx.NtRaiseHardError(
		_STATUS_INVALID_IMAGE_FORMAT,
		0,
		0,
		0,
		_OPTION_SHUTDOWN_SYSTEM,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}
