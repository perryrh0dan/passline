package storage

import "passline/pkg/crypt"

func DecryptCredential(credential *Credential, globalPassword []byte) error {
	err := decryptPassword(credential, globalPassword)
	if err != nil {
		return err
	}

	err = decryptRecoveryCodes(credential, globalPassword)
	if err != nil {
		return err
	}

	return nil
}

func EncryptCredential(credential *Credential, key []byte) error {
	err := encryptPassword(credential, key)
	if err != nil {
		return err
	}

	err = encryptRecoveryCodes(credential, key)
	if err != nil {
		return err
	}

	return nil
}

func encryptPassword(credential *Credential, key []byte) error {
	var err error
	credential.Password, err = crypt.AesGcmEncrypt(key, credential.Password)
	if err != nil {
		return err
	}

	return nil
}

func encryptRecoveryCodes(credential *Credential, globalPassword []byte) error {
	var encryptedRecoveryCodes = make([]string, 0)

	for _, c := range credential.RecoveryCodes {
		encryptedRecoveryCode, err := crypt.AesGcmEncrypt(globalPassword, c)
		if err != nil {
			return err
		}
		encryptedRecoveryCodes = append(encryptedRecoveryCodes, encryptedRecoveryCode)
	}

	credential.RecoveryCodes = encryptedRecoveryCodes
	return nil
}

func decryptPassword(credential *Credential, globalPassword []byte) error {
	// Decrypt passwords
	var err error
	credential.Password, err = crypt.AesGcmDecrypt(globalPassword, credential.Password)
	if err != nil {
		return err
	}

	return nil
}

func decryptRecoveryCodes(credential *Credential, globalPassword []byte) error {
	var decryptedRecoveryCodes = make([]string, 0)
	for _, c := range credential.RecoveryCodes {
		decryptedRecoveryCode, err := crypt.AesGcmDecrypt(globalPassword, c)
		if err != nil {
			return err
		}
		decryptedRecoveryCodes = append(decryptedRecoveryCodes, decryptedRecoveryCode)
	}

	credential.RecoveryCodes = decryptedRecoveryCodes
	return nil
}
