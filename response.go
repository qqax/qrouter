package router

import (
	"bytes"
	"encoding/json"
	"github.com/gabriel-vasile/mimetype"
	"github.com/rs/zerolog/log"
	"mime/multipart"
	"net/http"
)

func writeJSONResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		writeJSONError(w, "failed to write json", http.StatusInternalServerError)
		return
	}
}
func writeBytesResponse(w http.ResponseWriter, v []byte) {
	_, err := w.Write(v)
	if err != nil {
		writeJSONError(w, "failed to write response", http.StatusInternalServerError)
		return
	}

	m := mimetype.Detect(v)

	w.Header().Set("Content-Type", m.String())
	w.Header().Set("Accept-Ranges", "bytes")
}
func writeMultipartResponse(w http.ResponseWriter, v map[string][]byte) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for fieldName, field := range v {
		part, _ := writer.CreateFormField(fieldName)

		_, err := part.Write(field)
		if err != nil {
			writeJSONError(w, "failed to write form data", http.StatusInternalServerError)
			return
		}
	}

	err := writer.Close()
	if err != nil {
		writeJSONError(w, "failed to write form data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", writer.FormDataContentType())
}
func writeJSONError(writer http.ResponseWriter, text string, code int) {
	writer.WriteHeader(code)
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(map[string]string{"error": text})
	if err != nil {
		log.Error().Err(err).Msg("failed to write http error")
	}
}
func redirectJSONResponse(writer http.ResponseWriter, request *http.Request, url string, v any) {
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(v)
	if err != nil {
		writeJSONError(writer, "failed to write json", http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, url, 200)
}
