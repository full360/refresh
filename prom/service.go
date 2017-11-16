package prom

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"gitlab.full360.com/full360/refresh/storage"
)

type PromService interface {
	Refresh() (*http.Response, error)
}

type promService struct {
	storage    storage.S3Storage
	logger     log.Logger
	httpClient *http.Client
	config     struct {
		url    string
		method string
	}
}

func NewPromService(storage storage.S3Storage, logger log.Logger, client *http.Client, url, method string) *promService {
	return &promService{
		storage:    storage,
		logger:     logger,
		httpClient: client,
		config: struct {
			url    string
			method string
		}{
			url:    url,
			method: method,
		},
	}
}

func (p *promService) call() (*http.Response, error) {
	r, err := http.NewRequest(p.config.method, p.config.url, nil)
	if err != nil {
		return nil, err
	}

	return p.httpClient.Do(r)
}

func (p *promService) Refresh() (*http.Response, error) {
	if err := p.storage.Download(); err != nil {
		return nil, err
	}

	resp, err := p.call()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
