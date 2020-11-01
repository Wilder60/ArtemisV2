package test

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"testing"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/db"
)

const FirestoreEmulatorHost = "FIRESTORE_EMULATOR_HOST"

var testClient *db.Firestore
var procesID int

// Some test data used for unit testing
var testEvents = []domain.Event{
	domain.Event{
		ID:          "TESTID1",
		UserID:      "TESTUSERID",
		Type:        "EVENT",
		Name:        "TESTEVENT",
		Description: "DESCRIPTION",
		Color:       "#FFFFFF",
		SDate:       "2020-10-28 03:15:10-04:00",
		EDate:       "2020-10-28 04:15:10-04:00",
	},
	domain.Event{
		ID:          "TESTID2",
		UserID:      "TESTUSERID",
		Type:        "EVENT",
		Name:        "TESTEVENT",
		Description: "DESCRIPTION",
		Color:       "#FFFFFF",
		SDate:       "2020-10-28 04:15:10-04:00",
		EDate:       "2020-10-28 05:15:10-04:00",
	},
	domain.Event{
		ID:          "TESTID3",
		UserID:      "TESTUSERID",
		Type:        "EVENT",
		Name:        "TESTEVENT",
		Description: "DESCRIPTION",
		Color:       "#FFFFFF",
		SDate:       "2020-10-28 05:15:10-04:00",
		EDate:       "2020-10-28 06:15:10-04:00",
	},
}

func FirestoreSetUp() {
	cmd := exec.Command("gcloud", "beta", "emulators", "firestore", "start", "--host-port=localhost")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	procesID := cmd.Process.Pid
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		// reading it's output
		buf := make([]byte, 256, 256)
		for {
			n, err := stderr.Read(buf[:])
			if err != nil {
				// until it ends
				if err == io.EOF {
					break
				}
				log.Fatalf("reading stderr %v", err)
			}

			if n > 0 {
				d := string(buf[:n])

				// only required if we want to see the emulator output
				log.Printf("%s", d)

				// checking for the message that it's started
				if strings.Contains(d, "Dev App Server is now running") {
					wg.Done()
				}

				// and capturing the FIRESTORE_EMULATOR_HOST value to set
				pos := strings.Index(d, FirestoreEmulatorHost+"=")
				if pos > 0 {
					host := d[pos+len(FirestoreEmulatorHost)+1 : len(d)-1]
					os.Setenv(FirestoreEmulatorHost, host)
				}
			}
		}
	}()

	wg.Wait()
	newFirestoreTestClient()
	m.Run()
}

func FirestoreTearDown() {
	testClient.Client.Close()
	syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}

func newFirestoreTestClient() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "test")
	if err != nil {
		log.Fatalf("firebase.NewClient err: %v", err)
	}

	test := db.Firestore{
		Client:     client,
		Collection: client.Collection("events"),
	}
	testClient = &test
}

func insertMockData() {
	ctx := context.Background()

	for _, val := range testEvents {
		_, _, err := testClient.Collection.Add(ctx, &val)
		if err == iterator.Done {
			break
		} else if err != nil {
			panic(err.Error())
		}
	}
}

func removeMockData() {
	ctx := context.Background()
	iter := testClient.Collection.Documents(ctx)
	batch := testClient.Client.Batch()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			panic(err.Error())
		}

		batch.Delete(doc.Ref)

		_, err = batch.Commit(ctx)
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestGetEventsPaginatedSuccess(T *testing.T) {
	insertMockData()
	defer removeMockData()
	res, err := testClient.GetEventsPaginated(context.Background(), "TESTUSERID", "2020-10-28 03:15:10-04:00", 10, 0)
	if err != nil {
		T.Error(err)
	}
	if !cmp.Equal(res, testEvents) {
		T.Errorf("Invalid data returned from GetEventsPaginated\n Expected: %v\n Got: %v\n", testEvents, res)
	}

}

func TestGetEventsPaginatedEmpty(T *testing.T) {
	insertMockData()
	defer removeMockData()
	emptyEvents := []domain.Event{}
	res, err := testClient.GetEventsPaginated(context.Background(), "TESTUSERID", "9999-31-12 23:59:59-04:00", 10, 0)
	if err != nil {
		T.Error(err)
	}
	if !cmp.Equal(res, emptyEvents) {
		T.Errorf("Invalid data returned from GetEventsPaginated\n Expected: %v\n Got: %v\n", emptyEvents, res)
	}

}

func TestGetEventsInRange(T *testing.T) {
	insertMockData()
	defer removeMockData()
	res, err := testClient.GetEventsInRange(context.Background(), "TESTUSERID", "2020-10-28 03:15:10-04:00", "2020-10-28 06:15:10-04:00")
	if err != nil {
		T.Error(err)
	}

	if !cmp.Equal(res, testEvents) {
		T.Errorf("Invalid data returned from GetEventsInRange\n Expected: %v\n Got: %v\n", testEvents, res)
	}
}

func TestCreateEvents(T *testing.T) {
	defer removeMockData()
	ctx := context.Background()
	res := testClient.CreateEvents(ctx, testEvents[0])
	if res != nil {
		T.Errorf("Returned from CreateEvent request err: %v", res.Error())
	}

	iter := testClient.Collection.Where("ID", "==", testEvents[0].ID).Documents(ctx)
	doc, _ := iter.GetAll()
	if len(doc) != 1 {
		T.Errorf("Expected 1 document to return but retrieved %d", len(doc))
	}

	output := domain.Event{}
	err := doc[0].DataTo(&output)
	if err != nil {
		T.Error(err)
	}

	if !cmp.Equal(testEvents[0], output) {
		T.Errorf("Retrieved Documents do not match\n Expected: %v\n Got: %v\n", testEvents[0], output)
	}

}

func TestUpdateEvent(T *testing.T) {
	insertMockData()
	defer removeMockData()

	newEvent := domain.Event{
		ID:          "TESTID1",
		UserID:      "TESTUSERID",
		Type:        "NEWEVENT",
		Name:        "UPDATEDTESTEVENT",
		Description: "DESCRIPTION",
		Color:       "#FFFFFF",
		SDate:       "2020-10-28 03:15:10-04:00",
		EDate:       "2020-10-28 04:15:10-04:00",
	}

	ctx := context.Background()
	err := testClient.UpdateEvent(ctx, newEvent)
	if err != nil {
		T.Error(err)
	}
}

func TestDeleteEvents(T *testing.T) {
	insertMockData()
	defer removeMockData()

	ctx := context.Background()
	err := testClient.DeleteEvents(ctx, testEvents[0].UserID, testEvents[0].ID)
	if err != nil {
		T.Error(err)
	}
}

func TestDeleteEventsForUser(T *testing.T) {
	insertMockData()
	defer removeMockData()
	ctx := context.Background()
	expected := []domain.Event{}
	err := testClient.DeleteEventsForUser(ctx, "TESTUSERID")
	if err != nil {
		T.Error(err)
	}

	output := []domain.Event{}
	iter := testClient.Collection.Where("UserID", "==", "TESTUSERID").Documents(ctx)
	for {
		event := domain.Event{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			T.Error(err)
		}
		doc.DataTo(&event)
		output = append(output, event)
	}

	if !cmp.Equal(expected, output) {
		T.Errorf("Output does not match expected output\n Expected: %v\n Got: %v\n", expected, output)
	}

}
