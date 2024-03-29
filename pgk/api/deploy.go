package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Stolkerve/kappa/pgk/db"
	"github.com/Stolkerve/kappa/pgk/storage"
	"github.com/go-chi/chi/v5"
)

func DeployRoutes(r chi.Router) http.Handler {
	r.Post("/deploy", func(w http.ResponseWriter, r *http.Request) {
		functionWasmFile, functionWasmFileMetadata, err := r.FormFile("function")
		defer functionWasmFile.Close()
		if err != nil || functionWasmFileMetadata.Size > 1_048_576*5 || !strings.HasSuffix(functionWasmFileMetadata.Filename, ".wasm") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		functionWasmBuffer, err := io.ReadAll(functionWasmFile)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		funtion := storage.Function{
			Wasm: functionWasmBuffer,
		}

		if err := db.DB.Create(&funtion).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(funtion.ID))

	})

	r.Get("/deploys/{page:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.ParseUint(chi.URLParam(r, "page"), 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var ids []string
		if err := db.DB.Raw("SELECT id FROM functions lIMIT 50 OFFSET ?", page*50).Scan(&ids).Error; err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		jsonBuf, _ := json.Marshal(ids)
		w.Write(jsonBuf)
	})

	return r
}
