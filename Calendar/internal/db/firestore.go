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
	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain/requests"
	"go.uber.org/fx"
)

const collectionPath = "calendar/%s/events"

// Firestore is the sturct that contains a reference to the client and collection
type Firestore struct {
	Client *firestore.Client
}

// CreateFirestoreWrapper is the function that will be called to create the firestore
// to be used for fx's dependency injection
func CreateFirestoreWrapper(config *config.Config) *Firestore {

	ctx := context.Background()
	serviceAccount := option.WithCredentialsFile(config.Database.Firebase.ServiceAccount)
	app, err := firebase.NewApp(ctx, nil, serviceAccount)
	if err != nil {
		panic(fmt.Sprintf("Could not create new firebase app err: %v", err.Error()))
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		panic(fmt.Sprintf("Could not open firebase client err: %v", err.Error()))
	}

	return &Firestore{
		Client: client,
	}
}

// GetEventsPaginated will take a starting time, and will return all the events after
// sdate in a paginated format, using the limit and offset parameters
// GET api/v1/calendar?time=string&limit=int&offset=int
func (w Firestore) GetEventsPaginated(req *requests.GetPagination) ([]domain.Event, error) {
	var direction firestore.Direction
	if !req.Desc {
		direction = firestore.Asc
	} else {
		direction = firestore.Desc
	}

	collection := w.Client.Collection(fmt.Sprintf(collectionPath, req.UserID))
	iter := collection.OrderBy("SDate", direction).StartAt(req.Sdate).Offset(req.Offset).Limit(req.Limit).Documents(req.Ctx)
	return w.parseIterator(iter)
}

// GET api/v1/calendar/range?sdate=string&edate=string
func (w Firestore) GetEventsInRange(req *requests.GetRange) ([]domain.Event, error) {
	collection := w.Client.Collection(fmt.Sprintf(collectionPath, req.UserID))

	// Give me every event for a given user ordered by sdate, starting at the sdate and ending at edate
	iter := collection.OrderBy("SDate", firestore.Asc).StartAt(req.Sdate).EndAt(req.Edate).Documents(req.Ctx)
	return w.parseIterator(iter)
}

func (w Firestore) parseIterator(iter *firestore.DocumentIterator) ([]domain.Event, error) {
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
func (w Firestore) CreateEvents(req *requests.Add) error {
	collection := w.Client.Collection(fmt.Sprintf(collectionPath, req.UserID))
	event, err := req.ToEvent()
	if err != nil {
		return err
	}
	_, _, err = collection.Add(req.Ctx, &event)
	// log the add time
	if err != nil {
		return err
	}
	return nil
}

// PATCH api/v1/calendar
func (w Firestore) UpdateEvent(req *requests.Update) error {
	collection := w.Client.Collection(fmt.Sprintf(collectionPath, req.UserID))
	iterator := collection.Where("ID", "==", req.EventID).Documents(req.Ctx)
	// This should only return one document, but a check should be added just incase
	docs, err := iterator.GetAll()
	if err != nil {
		return err
	}

	if len(docs) == 0 {
		return errors.New("Error, no event found with matching ID")
	}
	_, err = docs[0].Ref.Set(req.Ctx, req.ToEvent())
	return err
}

// DELETE api/v1/calendar
func (w Firestore) DeleteEvents(req *requests.Delete) error {
	collection := w.Client.Collection(fmt.Sprintf(collectionPath, req.UserID))
	iter := collection.Where("ID", "==", req.ID).Documents(req.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		doc.Ref.Delete(req.Ctx)
	}
	return nil
}

// DELETE api/v1/calendar
func (w Firestore) DeleteEventsForUser(req *requests.Delete) error {

	var userDocPath string = "events/%s"
	collection := w.Client.Collection(fmt.Sprintf(collectionPath, req.UserID))

	var batchSize int = 100
	for {
		iter := collection.Where("UserID", "==", req.UserID).Limit(batchSize).Documents(req.Ctx)
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
			break
		}

		_, err := batch.Commit(req.Ctx)
		if err != nil {
			return err
		}
	}

	_, err := w.Client.Doc(fmt.Sprintf(userDocPath, req.UserID)).Delete(req.Ctx)
	return err
}

var FirebaseModule = fx.Option(
	fx.Provide(CreateFirestoreWrapper),
)
