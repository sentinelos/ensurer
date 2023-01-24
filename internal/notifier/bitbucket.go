package notifier

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/bitbucket"
	"github.com/drone/go-scm/scm/transport/oauth2"
)

// NewBitbucketNotifier returns a notifier that posts task steps statuses as status checks on a pull request.
func NewBitbucketNotifier(namespace, name, ref string) (*Bitbucket, error) {
	uri, ok := os.LookupEnv("BITBUCKET_URI")
	if !ok {
		uri = "https://api.bitbucket.org"
	}

	token, ok := os.LookupEnv("BITBUCKET_TOKEN")
	if !ok {
		return nil, errors.New("missing BITBUCKET_TOKEN")
	}

	client, err := bitbucket.New(uri)
	if err != nil {
		return nil, err
	}

	client.Client = &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&scm.Token{Token: token}),
		},
	}

	return &Bitbucket{
		GitNotifier: &GitNotifier{
			Namespace: namespace,
			Name:      name,
			Ref:       ref,
			client:    client,
			ctx:       context.Background(),
		},
	}, nil
}
