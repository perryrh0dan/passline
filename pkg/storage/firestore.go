package storage

import (
	"log"

	"golang.org/x/net/context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/perryrh0dan/passline/pkg/structs"
)

type FireStore struct {
	client *firestore.Client
}

func (fs *FireStore) Init() error {
	ctx := context.Background()

	opt := option.WithCredentialsFile("C:\\Users\\tpoe\\go\\src\\github.com\\perryrh0dan\\passline\\todo-83ef9-firebase-adminsdk-86yi4-73d003112b.json")
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

func (fs FireStore) GetByName(name string) (structs.Item, error) {
	dsnap, err := fs.client.Collection("passline").Doc(name).Get(context.Background())
	if err != nil {
		return structs.Item{}, err
	}
	var item structs.Item
	dsnap.DataTo(&item)

	return item, nil
}

func (fs FireStore) GetByIndex(index int) (structs.Item, error) {
	return structs.Item{}, nil
}

func (fs FireStore) GetAll() ([]structs.Item, error) {
	return []structs.Item{}, nil
}

func (fs FireStore) Add(item structs.Item) error {
	_, err := fs.client.Collection("passline").Doc(item.Name).Set(context.Background(), item)
	if err != nil {
		log.Fatalf("Failed adding aturing: %v", err)
	}
	return nil
}

func (fs FireStore) Delete(structs.Item) error {
	return nil
}
