package storage

import (
	"errors"
	"log"
	"os"
	"path"
	"sort"

	"golang.org/x/net/context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FireStore struct {
	client *firestore.Client
}

func NewFirestore() (*FireStore, error) {
	mainDir, _ := getMainDir()
	credentialsFile := path.Join(mainDir, "firestore.json")

	// Check for credentials file
	_, err := os.Stat(credentialsFile)
	if err != nil {
		log.Fatalf("error missing firebase credentials file\n")
		return nil, err
	}

	ctx := context.Background()

	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &FireStore{client: client}, nil
}

func (fs *FireStore) GetItemByName(name string) (Item, error) {
	dsnap, err := fs.client.Collection("passline").Doc(name).Get(context.Background())
	if err != nil {
		return Item{}, err
	}
	var item Item
	dsnap.DataTo(&item)

	return item, nil
}

func (fs *FireStore) GetItemByIndex(index int) (Item, error) {
	items, err := fs.GetAllItems()
	if err != nil {
		return Item{}, err
	}

	if index < 0 && index >= len(items) {
		return Item{}, errors.New("Out of index")
	}

	return items[index], nil
}

func (fs *FireStore) GetAllItems() ([]Item, error) {
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

	sort.Sort(ByName(items))
	return items, nil
}

func (fs *FireStore) CreateItem(item Item) error {
	_, err := fs.client.Collection("passline").Doc(item.Name).Set(context.Background(), item)
	if err != nil {
		log.Fatalf("Failed adding item: %v", err)
	}

	return nil
}

func (fs *FireStore) AddCredential(name string, credential Credential) error {
	item, err := fs.GetItemByName(name)
	if err != nil {
		return err
	}

	item.Credentials = append(item.Credentials, credential)
	err = fs.CreateItem(item)
	if err != nil {
		log.Fatalf("Failed updating credentials: %v", err)
	}

	return nil
}

func (fs *FireStore) deleteItem(item Item) error {
	_, err := fs.client.Collection("passline").Doc(item.Name).Delete(context.Background())
	if err != nil {
		log.Printf("An error has occured: %s", err)
		return err
	}

	return nil
}

func (fs *FireStore) DeleteCredential(item Item, credential Credential) error {
	indexCredential := getIndexOfCredential(item.Credentials, credential)
	if indexCredential == -1 {
		return errors.New("Item not found")
	}

	if len(item.Credentials) > 1 {
		item.Credentials = removeFromCredentials(item.Credentials, indexCredential)
		err := fs.CreateItem(item)
		if err != nil {
			log.Fatalf("Failed updating credentials: %v", err)
		}
	} else {
		fs.deleteItem(item)
	}

	return nil
}

func (fs *FireStore) UpdateItem(item Item) error {
	err := fs.deleteItem(item)
	if err != nil {
		return err
	}

	err = fs.CreateItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FireStore) GetAllItemNames() ([]string, error) {
	var names []string
	items, err := fs.GetAllItems()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		names = append(names, item.Name)
	}
	return names, nil
}
