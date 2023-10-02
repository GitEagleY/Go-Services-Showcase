package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// New initializes the 'client' variable with the provided MongoDB client instance and returns a new Models instance.
// This function is used to set up the MongoDB client for the Models struct.
func New(mongo *mongo.Client) Models {
	// Set the global 'client' variable to the provided MongoDB client instance.
	client = mongo
	// Return a Models instance with an empty LogEntry.
	return Models{LogEntry: LogEntry{}}
}

type Models struct {
	LogEntry LogEntry
}

// LogEntry is a struct representing a log entry document in MongoDB.
type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Insert inserts a new log entry into the MongoDB collection named 'logs'.
// It takes a LogEntry instance as input and returns an error if the insertion fails.
func (l *LogEntry) Insert(entry LogEntry) error {
	// Get the 'logs' collection from the MongoDB database.
	collection := client.Database("logs").Collection("logs")

	// Insert the provided LogEntry instance into the 'logs' collection.
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		ID:        entry.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}
	return nil
}

// All retrieves all log entries from the MongoDB collection 'logs' and returns them as a slice of LogEntry pointers.
// It also sorts the entries by the 'created_at' field in descending order.
func (l *LogEntry) All() ([]*LogEntry, error) {
	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get the 'logs' collection from the MongoDB database.
	collection := client.Database("logs").Collection("logs")

	// Define options for the find operation, including sorting by 'created_at' in descending order.
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	// Perform the find operation to retrieve all log entries.
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	// Create a slice to hold the retrieved log entries.
	var logs []*LogEntry

	// Iterate through the cursor and decode each entry into a LogEntry instance, appending it to the logs slice.
	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error decoding log into slice", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}
	return logs, nil
}

// GetOne retrieves a single log entry from the MongoDB collection 'logs' by its ID.
// It returns the retrieved log entry as a pointer to LogEntry and an error if the operation fails.
func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get the 'logs' collection from the MongoDB database.
	collection := client.Database("logs").Collection("logs")

	// Convert the provided ID string to an ObjectID.
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Define a LogEntry variable to store the retrieved log entry.
	var entry LogEntry

	// Find and decode the log entry by its ID.
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// DropCollection deletes the entire 'logs' collection from the MongoDB database.
// It returns an error if the collection drop operation fails.
func (l *LogEntry) DropCollection() error {
	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get the 'logs' collection from the MongoDB database.
	collection := client.Database("logs").Collection("logs")

	// Drop the 'logs' collection.
	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

// Update updates a log entry in the MongoDB collection 'logs' by its ID.
// It takes the updated fields from the LogEntry instance and returns the MongoDB UpdateResult and an error if the operation fails.
func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	// Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get the 'logs' collection from the MongoDB database.
	collection := client.Database("logs").Collection("logs")

	// Convert the LogEntry's ID to an ObjectID.
	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	// Perform an update operation to modify the specified log entry by its ID.
	result, err := collection.UpdateOne(ctx,
		bson.M{"id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		})
	if err != nil {
		return nil, err
	}
	return result, nil
}
