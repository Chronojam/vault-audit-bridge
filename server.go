package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net"

	"cloud.google.com/go/datastore"

	"golang.org/x/net/context"
)

var (
	googleCloudProject = flag.String("google.project", "", "Google cloud project")
	datastoreEntity    = flag.String("datastore.entity", "vault-audit", "Datastore Entity")
	addr               = flag.String("addr", "0.0.0.0:3333", "address to listen on")
	proto              = flag.String("proto", "tcp", "protocol listen on")
)

// AuditEntry represents a vault audit entry
type AuditEntry struct {
	Timestamp string
	Path      string
	Entry     string
}

// UnmarshalJSON implements JSON unsmarshaling for a vault audit event
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

// AuditHandler is a GCP DataStore based audit event handler
type AuditHandler struct {
	bqClient *datastore.Client
}

// HandleRequest incoming requests.
func (ah *AuditHandler) HandleRequest(conn net.Conn) {
	ctx := context.Background()
	// Close the connection when you're done with it.
	defer conn.Close()

	key := datastore.NameKey(*datastoreEntity, "", nil)
	dec := json.NewDecoder(conn)
	for {
		ae := &AuditEntry{}
		err := dec.Decode(ae)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Printf("error decoding, %v", err)
			continue
		}

		_, err = ah.bqClient.Put(ctx, key, ae)
		if err != nil {
			log.Printf("error adding to datastore, %v", err)
		}
	}
}

func main() {
	flag.Parse()
	ctx := context.Background()

	d, err := datastore.NewClient(ctx, *googleCloudProject)
	if err != nil {
		log.Fatalf("failed creating datastore client, %v", err)
	}

	ah := &AuditHandler{
		bqClient: d,
	}

	// Listen for incoming connections.
	l, err := net.Listen(*proto, *addr)
	if err != nil {
		log.Fatalf("error listening, %v", err)
	}
	// Close the listener when the application closes.
	defer l.Close()

	log.Printf("Listening on %s(%s)", *proto, *addr)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Printf("error accepting connection, %v", err)
			continue
		}
		// Handle connections in a new goroutine.
		go ah.HandleRequest(conn)
	}
}
