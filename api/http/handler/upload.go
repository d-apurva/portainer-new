package handler

import (
	"github.com/portainer/portainer"
	httperror "github.com/portainer/portainer/http/error"
	"github.com/portainer/portainer/http/security"

	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// UploadHandler represents an HTTP API handler for managing file uploads.
type UploadHandler struct {
	*mux.Router
	Logger      *log.Logger
	FileService portainer.FileService
}

// NewUploadHandler returns a new instance of UploadHandler.
func NewUploadHandler(bouncer *security.RequestBouncer) *UploadHandler {
	h := &UploadHandler{
		Router: mux.NewRouter(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	h.Handle("/upload/tls/{endpointID}/{certificate:(?:ca|cert|key)}",
		bouncer.AuthenticatedAccess(http.HandlerFunc(h.handlePostUploadTLS)))
	return h
}

func (handler *UploadHandler) handlePostUploadTLS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httperror.WriteMethodNotAllowedResponse(w, []string{http.MethodPost})
		return
	}

	vars := mux.Vars(r)
	endpointID := vars["endpointID"]
	certificate := vars["certificate"]
	ID, err := strconv.Atoi(endpointID)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	file, _, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}

	var fileType portainer.TLSFileType
	switch certificate {
	case "ca":
		fileType = portainer.TLSFileCA
	case "cert":
		fileType = portainer.TLSFileCert
	case "key":
		fileType = portainer.TLSFileKey
	default:
		httperror.WriteErrorResponse(w, portainer.ErrUndefinedTLSFileType, http.StatusInternalServerError, handler.Logger)
		return
	}

	err = handler.FileService.StoreTLSFile(portainer.EndpointID(ID), fileType, file)
	if err != nil {
		httperror.WriteErrorResponse(w, err, http.StatusInternalServerError, handler.Logger)
		return
	}
}
