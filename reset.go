package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only accessible in a development environments"))
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("couldn't reset the db"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted all users and Hits reset to 0"))
}
