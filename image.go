package qrouter

import (
	"context"
	"net/http"
	"strings"
)

type Imager interface {
	SetName(name string)
	ReadRawByName(ctx context.Context) error
	GetRaw() []byte
}

func GetImage(w http.ResponseWriter, r *http.Request, i *Imager) {
	path := strings.Split(r.URL.Path, "/")

	name := path[len(path)-1]
	if name == "" {
		WriteJSONError(w, "image id is empty", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	(*i).SetName(name)

	err := (*i).ReadRawByName(ctx)
	if err != nil {
		WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteBytesResponse(w, (*i).GetRaw())
}
