package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Stolkerve/kappa/pgk/db"
	"github.com/Stolkerve/kappa/pgk/storage"
	"github.com/Stolkerve/kappa/pgk/types"
	"github.com/Stolkerve/kappa/pgk/wasm"
	"github.com/go-chi/chi/v5"
)

func CallRoutes(r chi.Router) {
	r.HandleFunc("/call/{id}/*", func(w http.ResponseWriter, r *http.Request) {
		functionID := chi.URLParam(r, "id")

		var function storage.Function
		functionWasm := storage.Cache.Get(functionID)
		if functionWasm == nil {
			if db.DB.First(&function, "id = ?", functionID).Error != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Function not found :p"))
				return
			}
			storage.Cache.Push(function)
			functionWasm = function.Wasm
		}

		req, err := types.NewRequestWrapperFromHttpRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		runtimeRes, err := wasm.RunFunction(functionWasm, r.Context(), req)
		call := storage.Call{
			Stdout:      runtimeRes.Stdout,
			Stderr:      runtimeRes.Stderr,
			Duration:    runtimeRes.Duration,
			MemoryUsage: runtimeRes.Memory,
			Fail:        false,
			FunctionID:  functionID,
		}

		if err != nil {
			call.Fail = true
			call.ErrorMsg = err.Error()
			db.DB.Create(&call)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		db.DB.Create(&call)

		for k, v := range runtimeRes.ResponseWrapper.Header {
			for _, v := range v {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(runtimeRes.ResponseWrapper.StatusCode)
		w.Write(req.Body)

	})

	r.Get("/calls/{functionID}/{offset:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		offset, err := strconv.ParseUint(chi.URLParam(r, "offset"), 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var calls []storage.Call

		if err := db.DB.Where("function_id = ?", chi.URLParam(r, "functionID")).Find(&calls).Limit(50).Offset(int(offset)).Error; err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonBuf, _ := json.Marshal(calls)
		w.Write(jsonBuf)
	})
}
