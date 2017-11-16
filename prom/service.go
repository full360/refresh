package prom

import (
	"net/http"

	"gitlab.full360.com/full360/refresh/storage"
)

type Service interface {
	Refresh() (*http.Response, error)
}

type service struct {
	storage    storage.S3Storage
	httpClient *http.Client
	config     struct {
		url    string
		method string
	}
}

func NewService(storage storage.S3Storage, client *http.Client, url, method string) *service {
	return &service{
		storage:    storage,
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

func (p *service) call() (*http.Response, error) {
	r, err := http.NewRequest(p.config.method, p.config.url, nil)
	if err != nil {
		return nil, err
	}

	return p.httpClient.Do(r)
}

func (p *service) Refresh() (*http.Response, error) {
	if err := p.storage.Download(); err != nil {
		return nil, err
	}

	resp, err := p.call()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
