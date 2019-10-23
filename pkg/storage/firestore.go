package storage

import (
	"errors"
	"log"
	"path"

	"golang.org/x/net/context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FireStore struct {
	client *firestore.Client
}

func (fs *FireStore) Init() error {
	mainDir, _ := getMainDir()
	credentialsFile := path.Join(mainDir, "firebase.json")

	ctx := context.Background()

	opt := option.WithCredentialsFile(credentialsFile)
	config := &firebase.Config{ProjectID: "todo-83ef9"}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	fs.client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func (fs FireStore) GetByName(name string) (Item, error) {
	dsnap, err := fs.client.Collection("passline").Doc(name).Get(context.Background())
	if err != nil {
		return Item{}, err
	}
	var item Item
	dsnap.DataTo(&item)

	return item, nil
}

func (fs FireStore) GetByIndex(index int) (Item, error) {
	items, err := fs.GetAll()
	if err != nil {
		return Item{}, err
	}

	if index < 0 && index >= len(items) {
		return Item{}, errors.New("Out of index")
	}

	return items[index], nil
}

func (fs FireStore) GetAll() ([]Item, error) {
	items := []Item{}
	iter := fs.client.Collection("passline").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var item Item
		doc.DataTo(&item)
		items = append(items, item)
	}

	return items, nil
}

func (fs FireStore) AddItem(item Item) error {
	_, err := fs.client.Collection("passline").Doc(item.Name).Set(context.Background(), item)
	if err != nil {
		log.Fatalf("Failed adding item: %v", err)
	}

	return nil
}

func (fs FireStore) AddCredential(name string, credential Credential) error {
	item, err := fs.GetByName(name)
	if err != nil {
		return err
	}

	item.Credentials = append(item.Credentials, credential)
	err = fs.AddItem(item)
	if err != nil {
		log.Fatalf("Failed updating credentials: %v", err)
	}

	return nil
}

func (fs FireStore) DeleteItem(item Item) error {
	_, err := fs.client.Collection("passline").Doc(item.Name).Delete(context.Background())
	if err != nil {
		log.Printf("An error has occured: %s", err)
		return err
	}

	return nil
}

func (fs FireStore) DeleteCredential(item Item, credential Credential) error {
	indexCredential := getIndexOfCredential(item.Credentials, credential)
	if indexCredential == -1 {
		return errors.New("Item not found")
	}

	item.Credentials = removeFromCredentials(item.Credentials, indexCredential)
	err := fs.AddItem(item)
	if err != nil {
		log.Fatalf("Failed updating credentials: %v", err)
	}

	return nil
}
