package numan

//API exposes overall application interface
type API interface {
	NumberAPI
	HistoryAPI
	Close()
}
