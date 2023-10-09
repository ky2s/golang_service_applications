package objects

type URLResponseData struct {
	ShortURL string `json:"short_url"`
	URL      string `json:"url"`
}
type ShortURLResponse struct {
	Data   URLResponseData `json:"data"`
	Status int             `json:"status"`
}
