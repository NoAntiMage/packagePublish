package request

type HttpClient interface {
}

func NewHttpClient() {}

type httpClient struct {
	token string
}
