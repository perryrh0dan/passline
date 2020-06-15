package main

import (
	"context"
	"passline/pkg/config"
	"passline/pkg/crypt"
	"passline/pkg/storage"
)

const (
	oldPassword = ""
	newPassword = ""
)

// 0.7.3 -> 0.8.0
func migrateV1() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	store, err := storage.New(cfg)
	if err != nil {
		return err
	}

	items, err := store.GetAllItems(context.TODO())
	if err != nil {
		return err
	}

	// get pwKey
	pwKey := crypt.GetKey([]byte(oldPassword))

	// decrypt old data
	for i := 0; i < len(items); i++ {
		for x := 0; x < len(items[i].Credentials); x++ {
			plainPW, err := crypt.AesGcmDecrypt([]byte(pwKey), items[i].Credentials[x].Password)
			items[i].Credentials[x].Password = plainPW
			if err != nil {
				return err
			}
		}
	}

	// encrypt with new logic
	encryptionKey, err := crypt.GenerateKey()
	if err != nil {
		return err
	}

	for i := 0; i < len(items); i++ {
		for x := 0; x < len(items[i].Credentials); x++ {
			pw, err := crypt.AesGcmEncrypt([]byte(encryptionKey), items[i].Credentials[x].Password)
			items[i].Credentials[x].Password = pw
			if err != nil {
				return err
			}
		}
	}

	// encrypt encryption key
	pwKey = crypt.GetKey([]byte(newPassword))
	encryptedEncryptionKey, err := crypt.AesGcmEncrypt(pwKey, encryptionKey)
	if err != nil {
		return err
	}

	data := storage.Data{
		Key:   encryptedEncryptionKey,
		Items: items,
	}

	err = store.SetData(context.TODO(), data)
	if err != nil {
		return err
	}

	return nil
}
