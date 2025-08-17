package httpservice

import (
	"encoding/json"
	"net/http"

	"github.com/shamhub/exercise/pkg/errorlib"
)

func SendErrorResponse(w http.ResponseWriter, c *customError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c.GetStatusCode())
	encodedError, err := json.Marshal(c)
	if err != nil {
		marshalErrorResponse(w, err)
	}
	w.Write(encodedError)
}

// marshalErrorResponse() should be invoked
// if json.Marshal() fails
func marshalErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	response := errorlib.MarshalErrorResponse{
		Message: err.Error(),
	}
	json.NewEncoder(w).Encode(response)
}
