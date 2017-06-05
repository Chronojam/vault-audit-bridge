package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"cloud.google.com/go/datastore"

	"golang.org/x/net/context"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var (
	GoogleCloudProject = flag.String("google.project", "", "Google cloud project")
	DatastoreEntity    = flag.String("datastore.entity", "vault-audit", "Datastore Entity")
)

type AuditEntry struct {
	Timestamp string
	Path      string
	Entry     string
}

func (a *AuditEntry) UnmarshalJSON(data []byte) error {
	p := map[string]interface{}{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	req := p["request"].(map[string]interface{})
	a.Path = req["path"].(string)
	a.Timestamp = p["time"].(string)
	a.Entry = string(data)

	return nil
}

type AuditHandler struct {
	bqClient *datastore.Client
}

func main() {
	flag.Parse()
	ctx := context.Background()
	d, err := datastore.NewClient(ctx, *GoogleCloudProject)
	if err != nil {
		log.Fatal(err.Error())
	}

	ah := &AuditHandler{
		bqClient: d,
	}

	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go ah.HandleRequest(conn)
	}
}

// Handles incoming requests.
func (ah *AuditHandler) HandleRequest(conn net.Conn) {
	ctx := context.Background()
	// Close the connection when you're done with it.
	defer conn.Close()

	key := datastore.NameKey(*DatastoreEntity, "", nil)
	dec := json.NewDecoder(conn)
	for {
		ae := &AuditEntry{}
		err := dec.Decode(ae)
		if err != nil {
			log.Printf("Error Decoding", err.Error())
			continue
		}

		_, err = ah.bqClient.Put(ctx, key, ae)
		if err != nil {
			log.Printf(err.Error())
		}
	}

}
