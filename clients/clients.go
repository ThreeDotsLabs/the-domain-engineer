package clients

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ThreeDotsLabs/the-domain-engineer/clients/bank"
	"github.com/ThreeDotsLabs/the-domain-engineer/clients/files"
	"github.com/ThreeDotsLabs/the-domain-engineer/clients/tax"
)

type Clients struct {
	Bank  bank.ClientWithResponsesInterface
	Files files.ClientWithResponsesInterface
	Tax   tax.ClientWithResponsesInterface
}

func NewClients(
	gatewayAddress string,
	requestEditorFn RequestEditorFn,
) (*Clients, error) {
	return NewClientsWithHttpClient(gatewayAddress, requestEditorFn, http.DefaultClient)
}

func NewClientsWithHttpClient(
	gatewayAddress string,
	requestEditorFn RequestEditorFn,
	httpDoer HttpDoer,
) (*Clients, error) {
	if gatewayAddress == "" {
		return nil, fmt.Errorf("gateway address is required")
	}

	if requestEditorFn == nil {
		requestEditorFn = func(_ context.Context, _ *http.Request) error {
			return nil
		}
	}

	bankClient, err := newClient(
		gatewayAddress,
		"bank-api",
		bank.NewClientWithResponses,
		bank.WithRequestEditorFn(bank.RequestEditorFn(requestEditorFn)),
		bank.WithHTTPClient(httpDoer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create bank client: %w", err)
	}

	filesClient, err := newClient(
		gatewayAddress,
		"files-api",
		files.NewClientWithResponses,
		files.WithRequestEditorFn(files.RequestEditorFn(requestEditorFn)),
		files.WithRequestEditorFn(files.RequestEditorFn(requestEditorFn)),
		files.WithHTTPClient(httpDoer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create files client: %w", err)
	}

	taxClient, err := newClient(
		gatewayAddress,
		"tax-api",
		tax.NewClientWithResponses,
		tax.WithRequestEditorFn(tax.RequestEditorFn(requestEditorFn)),
		tax.WithHTTPClient(httpDoer),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tax client: %w", err)
	}

	return &Clients{
		Bank:  bankClient,
		Files: filesClient,
		Tax:   taxClient,
	}, nil
}

type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestEditorFn func(ctx context.Context, req *http.Request) error

func newClient[Client any, ClientOption any](
	gatewayAddress string,
	serviceName string,
	clientConstructor func(server string, opts ...ClientOption) (*Client, error),
	requestEditorFn ...ClientOption,
) (*Client, error) {
	apiServerAddr, err := url.JoinPath(gatewayAddress, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create files api server address of %s: %w", serviceName, err)
	}

	apiClient, err := clientConstructor(apiServerAddr, requestEditorFn...)
	if err != nil {
		return nil, fmt.Errorf("failed to create files client: %w", err)
	}

	return apiClient, nil
}
