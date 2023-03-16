package arango

import (
	"context"

	arangodriver "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

type Client struct {
	DriverConnection arangodriver.Connection
	DriverClient     arangodriver.Client
}

func NewClient(addrs []string) (*Client, error) {
	conn, err := arangohttp.NewConnection(arangohttp.ConnectionConfig{
		Endpoints: addrs,
	})
	if err != nil {
		return &Client{}, err
	}

	cl, err := arangodriver.NewClient(arangodriver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		return &Client{}, err
	}
	return &Client{DriverClient: cl}, nil
}

func (cl *Client) Database(ctx context.Context, name string, createOnNotExist bool, collections []string) (arangodriver.Database, error) {
	var err error
	dbClient, err := cl.DriverClient.Database(ctx, name)
	if err != nil {
		if arangodriver.IsNotFoundGeneral(err) && createOnNotExist {
			// Create database with name
			dbClient, err := cl.DriverClient.CreateDatabase(ctx, name, nil)
			if err != nil {
				return nil, err
			}
			// Create collections
			for _, colName := range collections {
				_, err = dbClient.CreateCollection(ctx, colName, nil)
				if err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}
	return dbClient, nil
}
