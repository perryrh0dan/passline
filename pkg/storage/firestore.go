package storage

import (
	"context"
	"errors"
	"log"
	"os"
	"passline/pkg/config"
	"path"
	"sort"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FireStore struct {
	client *firestore.Client
	items  []Item
}

const (
	DataCollection   = "passline"
	ConfigCollection = "config"
	DefaultCategory  = "default"
)

func NewFirestore() (*FireStore, error) {
	mainDir := config.Directory()
	credentialsFile := path.Join(mainDir, "firestore.json")

	ctx := context.TODO()

	// Check for credentials file
	_, err := os.Stat(credentialsFile)
	if err != nil {
		log.Fatalf("error missing firebase credentials file\n")
		return nil, err
	}

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

func (fs *FireStore) GetItemByName(ctx context.Context, name string) (Item, error) {
	dsnap, err := fs.client.Collection(DataCollection).Doc(name).Get(context.Background())
	if err != nil {
		return Item{}, err
	}
	var item Item
	dsnap.DataTo(&item)

	// Add default Category if not exists
	for index, cred := range item.Credentials {
		if cred.Category == "" {
			item.Credentials[index].Category = DefaultCategory
		}
	}

	return item, nil
}

func (fs *FireStore) GetItemByIndex(ctx context.Context, index int) (Item, error) {
	items, err := fs.GetAllItems(ctx)
	if err != nil {
		return Item{}, err
	}

	if index < 0 && index >= len(items) {
		return Item{}, errors.New("Out of index")
	}

	return items[index], nil
}

func (fs *FireStore) GetAllItems(ctx context.Context) ([]Item, error) {
	if len(fs.items) > 0 {
		return fs.items, nil
	}

	iter := fs.client.Collection(DataCollection).Documents(ctx)
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

		// Add default Category if not exists
		for index, cred := range item.Credentials {
			if cred.Category == "" {
				item.Credentials[index].Category = DefaultCategory
			}
		}

		fs.items = append(fs.items, item)
	}

	sort.Sort(ByName(fs.items))
	return fs.items, nil
}

func (fs *FireStore) AddCredential(ctx context.Context, name string, credential Credential) error {
	item, err := fs.GetItemByName(ctx, name)
	if status.Code(err) == codes.NotFound {
		item = Item{Name: name, Credentials: []Credential{credential}}
	} else if err != nil {
		return err
	} else {
		item.Credentials = append(item.Credentials, credential)
	}

	err = fs.createItem(ctx, item)
	if err != nil {
		log.Fatalf("Failed updating credentials: %v", err)
	}

	return nil
}

func (fs *FireStore) DeleteCredential(ctx context.Context, item Item, username string) error {
	indexCredential := getIndexOfCredential(item.Credentials, username)
	if indexCredential == -1 {
		return errors.New("Item not found")
	}

	if len(item.Credentials) > 1 {
		item.Credentials = removeFromCredentials(item.Credentials, indexCredential)
		err := fs.createItem(ctx, item)
		if err != nil {
			log.Fatalf("Failed updating credentials: %v", err)
		}
	} else {
		fs.deleteItem(ctx, item)
	}

	return nil
}

func (fs *FireStore) UpdateItem(ctx context.Context, item Item) error {
	err := fs.deleteItem(ctx, item)
	if err != nil {
		return err
	}

	err = fs.createItem(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FireStore) SetData(ctx context.Context, data Data) error {
	fs.deleteCollection(ctx, 100)
	batch := fs.client.Batch()

	for _, item := range data.Items {
		itemRef := fs.client.Collection(DataCollection).Doc(item.Name)
		batch.Set(itemRef, item)
	}

	itemRef := fs.client.Collection(ConfigCollection).Doc("config")
	batch.Set(itemRef, Config{Key: data.Key})

	_, err := batch.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (fs *FireStore) GetKey(ctx context.Context) (string, error) {
	doc, err := fs.client.Collection(ConfigCollection).Doc("config").Get(ctx)
	if err != nil {
		return "", nil
	}

	var config Config
	doc.DataTo(&config)
	return config.Key, nil
}

func (fs *FireStore) SetKey(ctx context.Context, key string) error {
	doc := fs.client.Collection(ConfigCollection).Doc("config")
	_, err := doc.Set(ctx, Config{Key: key})
	if err != nil {
		return err
	}

	return nil
}

func (fs *FireStore) createItem(ctx context.Context, item Item) error {
	_, err := fs.client.Collection(DataCollection).Doc(item.Name).Set(ctx, item)
	if err != nil {
		log.Fatalf("Failed adding item: %v", err)
	}

	return nil
}

func (fs *FireStore) deleteItem(ctx context.Context, item Item) error {
	_, err := fs.client.Collection(DataCollection).Doc(item.Name).Delete(ctx)
	if err != nil {
		log.Printf("An error has occured: %s", err)
		return err
	}

	return nil
}

func (fs *FireStore) deleteCollection(ctx context.Context, batchSize int) error {
	ref := fs.client.Collection("passline")

	for {
		// Get a batch of documents
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to a
		// WriteBatch.
		batch := fs.client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}
