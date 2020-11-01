package db

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
	"github.com/Wilder60/ArtemisV2/Calendar/config"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain"
	"go.uber.org/fx"
)

// Firestore is the sturct that contains a reference to the client and collection
type Firestore struct {
	Client     *firestore.Client
	Collection *firestore.CollectionRef
}

// CreateFirestoreWrapper is the function that will be called to create the firestore
// to be used for fx's dependency injection
func CreateFirestoreWrapper(config *config.Config) *Firestore {
	ctx := context.Background()
	serviceAccount := option.WithCredentialsFile(config.Database.Postgres.User)
	app, err := firebase.NewApp(ctx, nil, serviceAccount)
	if err != nil {
		panic(fmt.Sprintf("Could not create new firebase app err: %v", err.Error()))
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		panic(fmt.Sprintf("Could not open firebase client err: %v", err.Error()))
	}

	return &Firestore{
		Client:     client,
		Collection: client.Collection("events"),
	}
}

// GetEventsPaginated will take a starting time, and will return all the events after
// sdate in a paginated format, using the limit and offset parameters
// GET api/v1/calendar?time=string&limit=int&offset=int
func (w *Firestore) GetEventsPaginated(ctx context.Context, userID, sdate string, limit, offset int) ([]domain.Event, error) {
	iter := w.Collection.Where("UserID", "==", userID).
		OrderBy("SDate", firestore.Asc).
		StartAt(sdate).
		Offset(offset).Limit(limit).Documents(ctx)

	return w.parseIterator(iter)
}

// GET api/v1/calendar/range?sdate=string&edate=string
func (w *Firestore) GetEventsInRange(ctx context.Context, userID, sdate, edate string) ([]domain.Event, error) {
	// Give me every event for a given user ordered by sdate, starting at the sdate and ending at edate
	iter := w.Collection.Where("UserID", "==", userID).
		OrderBy("SDate", firestore.Asc).
		StartAt(sdate).
		EndAt(edate).
		Documents(ctx)

	return w.parseIterator(iter)
}

func (w *Firestore) parseIterator(iter *firestore.DocumentIterator) ([]domain.Event, error) {
	serializedEvents := []domain.Event{}
	for {
		var event = domain.Event{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}
		err = doc.DataTo(&event)
		if err != nil {
			fmt.Printf("Error serializing document to struct err: %v\n", err.Error())
		} else {
			serializedEvents = append(serializedEvents, event)
		}
	}
	return serializedEvents, nil
}

// POST api/v1/calendar
func (w *Firestore) CreateEvents(ctx context.Context, event domain.Event) error {
	_, _, err := w.Collection.Add(ctx, &event)
	// log the add time
	if err != nil {
		fmt.Printf("Create event failed err: %v\n", err.Error())
		return err
	}
	return nil
}

// PATCH api/v1/calendar
func (w *Firestore) UpdateEvent(ctx context.Context, event domain.Event) error {
	iterator := w.Collection.Where("UserID", "==", event.UserID).Where("ID", "==", event.ID).Documents(ctx)
	docs, err := iterator.GetAll()
	if err != nil {
		return err
	}
	if len(docs) == 0 {
		return errors.New("Error, no event found with matching ID")
	}
	_, err = docs[0].Ref.Set(ctx, &event)
	return err
}

// DELETE api/v1/calendar
func (w *Firestore) DeleteEvents(ctx context.Context, userID, eventID string) error {
	iter := w.Collection.Where("UserID", "==", userID).Where("ID", "==", eventID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		doc.Ref.Delete(ctx)
	}
	return nil
}

// DELETE api/v1/calendar
func (w *Firestore) DeleteEventsForUser(ctx context.Context, userID string) error {
	var batchSize int = 100
	for {
		iter := w.Collection.Where("UserID", "==", userID).Limit(batchSize).Documents(ctx)
		numDeleted := 0

		batch := w.Client.Batch()
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

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}

var FirebaseModule = fx.Option(
	fx.Provide(CreateFirestoreWrapper),
)
