package api

import (
	"encoding/json"
	"net/http"
)

type hostResponse struct {
	Hosts []hostItem `json:"hosts"`
}

type hostItem struct { //bir hostun apicevabında görünecek hali
	HostID   string `json:"host_id"`
	Status   string `json:"status"`
	LastSeen string `json:"last_seen"`
}

func (s *Server) handleListHosts(w http.ResponseWriter, r *http.Request) {
	hosts, err := s.store.ListHosts()
	if err != nil {
		http.Error(w, "host listesi alinamadi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := hostResponse{Hosts: make([]hostItem, 0, len(hosts))} //hostResponse structının Hosts alanını hosts slice'ının uzunluğu kadar önceden ayırır
	for _, h := range hosts {
		resp.Hosts = append(resp.Hosts, hostItem{
			HostID:   h.HostID,
			Status:   h.Status,
			LastSeen: h.LastSeen,
		})
	}

	w.Header().Set("Content-Type", "application/json") //cevabın json olduğu dashboarda bildiriliyor
	json.NewEncoder(w).Encode(resp)
}
