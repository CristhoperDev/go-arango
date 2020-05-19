package connection

import (
	"container/list"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	arangodb "github.com/arangodb/go-driver"
	arangodbhttp "github.com/arangodb/go-driver/http"
	"github.com/cristhoperdev/events-import/model"
	"io/ioutil"
	"os"
	"time"
)

// Datastore represents the embedded database
type Datastore struct {
	client    arangodb.Client
	databases map[string]*Db
}

// Db represents a database
type Db struct {
	Name        string
	db          arangodb.Database
	collections map[string]arangodb.Collection
}

// Concept represents a concept definition
type Concept struct {
	Name     string                   `json:"_key" reindex:"_key,hash,pk"`
	Hash     uint64                   `json:"_hash" hash:"ignore" reindex:"_hash"`
	IsShared bool                     `json:"isShared" reindex:"isShared"`
	Domain   string                   `json:"domain" reindex:"domain"`
	Cached   bool                     `json:"cached" reindex:"cached"`
}

// Open opens the datastore
func (s *Datastore) Open() error {

	var err error

	endpoints := []string{"http://localhost:8529"}

	login := ""
	pwd := ""

	conn, err := arangodbhttp.NewConnection(arangodbhttp.ConnectionConfig{
		Endpoints: endpoints,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})

	if err != nil {
		return err
	}

	client, err := arangodb.NewClient(arangodb.ClientConfig{
		Connection:     conn,
		Authentication: arangodb.BasicAuthentication(login, pwd),
	})


	if s.databases == nil {
		s.databases = make(map[string]*Db)
	}

	if err != nil {
		return err
	}

	s.client = client

	return nil
}

// GetDatabase returns the database with the specified name
func (s *Datastore) GetDatabase(name string) (*Db, error) {

	//We first try to retrieve it from the map
	result := s.databases[name]
	var err error

	// If found, we return it
	if result != nil {
		return result, nil
	}

	// We create or retrieve the database
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	found, err := s.client.DatabaseExists(ctx, name)

	var db arangodb.Database

	if err != nil {
		return nil, err
	}

	if !found {

		ctx := arangodb.WithWaitForSync(ctx)
		db, err = s.client.CreateDatabase(ctx, name, nil)

		if err != nil {
			return nil, err
		}

	} else {

		db, err = s.client.Database(ctx, name)

		if err != nil {
			return nil, err
		}
	}

	result = &Db{
		Name:        name,
		db:          db,
		collections: make(map[string]arangodb.Collection),
	}

	// We populate the collections
	result.ensureCollections()

	s.databases[name] = result

	return result, nil
}

func (s *Db) ensureCollections() error {
	var err error

	// We ensure the collection in ArangoDB
	err = s.ensureCollection("Event", s.EventIsEdgeCollection())
	if err != nil {
		return err
	}

	// We ensure the collection in ArangoDB
	err = s.ensureCollection("Movie", s.EventIsEdgeCollection())
	if err != nil {
		return err
	}

	return nil
}

// AnnouncementIsEdgeCollection defines if the collection is an edge collection
func (s *Db) EventIsEdgeCollection() bool {
	return false
}

// MovieIsEdgeCollection defines if the collection is an edge collection
func (s *Db) MovieIsEdgeCollection() bool {
	return false
}

func (s *Db) ensureCollection(collectionName string, isRelation bool) error {

	col, err := s.GetDocumentCollection(collectionName, isRelation)

	if err != nil {
		return err
	}

	s.collections[collectionName] = col

	return nil
}

// GetDocumentCollection creates or/and returns the specified document collection
func (s *Db) GetDocumentCollection(collection string, isRelation bool) (arangodb.Collection, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	ctx = arangodb.WithWaitForSync(ctx)

	found, err := s.db.CollectionExists(ctx, collection)

	if err != nil {

		return nil, err
	}

	if !found {

		var options *arangodb.CreateCollectionOptions

		if isRelation {
			options = &arangodb.CreateCollectionOptions{Type: arangodb.CollectionTypeEdge}
		} else {
			options = &arangodb.CreateCollectionOptions{Type: arangodb.CollectionTypeDocument}
		}

		col, err := s.db.CreateCollection(ctx, collection, options)
		col.EnsureTTLIndex(ctx, "createdAt", 60, nil)

		if err != nil {
			return nil, err
		}

		return col, nil
	}

	col, err := s.db.Collection(ctx, collection)

	if err != nil {
		return nil, err
	}

	return col, nil
}

// Seed populate the store with predefined data
func (s *Datastore) Seed() error {

	root := "seed"

	if _, err := os.Stat(root); !os.IsNotExist(err) {

		folders, err := ioutil.ReadDir(root)

		if err != nil {
			return err
		}

		// Each folder represents a database
		for _, folder := range folders {
			if folder.IsDir() {

				name := folder.Name()

				_, err := s.GetDatabase(name)

				if err != nil {
					return err
				}

				/*filepath.Walk(name, func(path string, info os.FileInfo, err error) error {

					if info.IsDir() == false && filepath.Ext(path) == ".json" {

						switch strings.ToLower(info.Name()) {
						case "announcements.json":
							return db.SeedAnnouncements(path)
						default:
							return nil
						}
					}

					return nil
				})*/
			}
		}
	}

	return nil
}


// Bulk import array in the store
func (s *Db) BulkImportEvents(item []*model.Event) error {

	_, err := s.ImportDBDocument(s.collections["Event"], item)

	if err != nil {
		return err
	}

	return nil
}


// Bulk import array in the store
func (s *Db) BulkImportEventLists(item *list.List) error {

	_, err := s.ImportDBDocument(s.collections["Event"], item)

	if err != nil {
		return err
	}

	return nil
}

//import document
func (s *Db) ImportDBDocument(collection arangodb.Collection, document interface{}) (arangodb.ImportDocumentStatistics, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	ctx = arangodb.WithWaitForSync(ctx)

	meta, err := collection.ImportDocuments(ctx, document, nil)

	if err != nil {
		return meta, err
	}

	return meta, nil
}

// CreateDBDocument creates the specified document within the passed collection
func (s *Db) CreateDBDocument(collection arangodb.Collection, document interface{}) (string, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	ctx = arangodb.WithWaitForSync(ctx)

	meta, err := collection.CreateDocument(ctx, document)

	if err != nil {
		return "", err
	}

	return meta.Key, nil
}

// UpdateDBDocument updates the specified document within the passed collection
func (s *Db) UpdateDBDocument(collection arangodb.Collection, documentID string, document interface{}) error {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	ctx = arangodb.WithWaitForSync(ctx)

	_, err := collection.UpdateDocument(ctx, documentID, document)

	if err != nil {
		return err
	}

	return nil
}

// DeleteDBDocument deletes the document with the specified id from the passed collection
func (s *Db) DeleteDBDocument(collection arangodb.Collection, documentID string) error {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	ctx = arangodb.WithWaitForSync(ctx)

	_, err := collection.RemoveDocument(ctx, documentID)

	if err != nil {
		return err
	}

	return nil
}

// GetAllMovies returns all the movies from the datastore
func (s *Db) GetAllMovies() ([]*model.Movie, error) {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*model.Movie

	var query string

	query = fmt.Sprintf("FOR d IN %s RETURN d", "Movie")

	cursor, err := s.db.Query(ctx, query, nil)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var doc model.Movie
		_, err := cursor.ReadDocument(ctx, &doc)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &doc)
	}

	return result, nil
}

// CreateEvent creates a event in the store
func (s *Db) CreateEvent(item *model.EventLog) error {

	_, err := s.CreateDBDocument(s.collections["Event"], item)

	if err != nil {
		return err
	}

	return nil
}

// CreateAnnouncement creates a announcement in the store
func (s *Db) CreateMovie(item *model.Movie) error {

	_, err := s.CreateDBDocument(s.collections["Movie"], item)

	if err != nil {
		return err
	}

	return nil
}

// UpdateAnnouncement updates the passed announcement in the store
func (s *Db) UpdateMovie(item *model.Movie) error {

	err := s.UpdateDBDocument(s.collections["Movie"], *item.Key, item)

	if err != nil {
		return err
	}

	return nil
}

// DeleteAnnouncement deletes the passed announcement from the datastore
func (s *Db) DeleteMovie(item *model.Movie) error {

	var id string

	id = *item.Key

	err := s.DeleteDBDocument(s.collections["Movie"], id)
	if err != nil {
		return err
	}

	return nil
}


// CreateEvent creates a event in the store
func (s *Db) UpdateDataEvent(key string, event []*model.Event) error {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	events, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	var query string

	query = "FOR d IN Event FILTER d._key == '" + key + "' UPDATE d WITH { events: APPEND(d.events," + string(events) + " ) } IN Event"

	cursor, err := s.db.Query(ctx, query, nil)
	if err != nil {
		return err
	}

	defer cursor.Close()

	return nil

	/*ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var result []*model.Movie

	var query string

	query = fmt.Sprintf("FOR d IN %s RETURN d", "Movie")

	cursor, err := s.db.Query(ctx, query, nil)

	if err != nil {
		return nil, err
	}

	defer cursor.Close()

	for {

		var doc model.Movie
		_, err := cursor.ReadDocument(ctx, &doc)

		if arangodb.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}

		result = append(result, &doc)
	}

	return result, nil*/

	return nil
}
