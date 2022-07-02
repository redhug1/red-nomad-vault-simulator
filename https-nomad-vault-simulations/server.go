package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/nomad/structs"

	"github.com/go-chi/chi/v5"
	//	"github.com/go-chi/chi/v5/middleware"
)

const (
	CertPath string = "../tls-certificates/00-certificates/server/cert.pem"
	KeyPath  string = "../tls-certificates/00-certificates/server/key.pem"
)

func main() {
	// Nomad is handled via normal go 'NewServeMux'
	mux := http.NewServeMux()

	// add an endpoint to show this is an HTTPS server with certificates
	mux.HandleFunc("/server", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "i am protected")
	})

	// Nomad 'healthcheck' simulation
	mux.HandleFunc("/v1/agent/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Nomad mock: health")

		w.Header().Set("Server", "A Go Nomad Mock Server")
		w.WriteHeader(200)
	})

	// Nomad 'plan' simulation
	mux.HandleFunc("/v1/job/dp-search-data-finder/plan", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Nomad mock: got plan\n")

		time.Sleep(100 * time.Millisecond) // small delay to simulate work being done

		w.Header().Set("Content-Type", "application/json")

		var res api.JobPlanResponse

		jsonResp, err := json.Marshal(res)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	})

	// Nomad 'run' simulation
	mux.HandleFunc("/v1/jobs", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Nomad mock: got run\n")

		time.Sleep(100 * time.Millisecond) // small delay to simulate work being done

		w.Header().Set("Content-Type", "application/json")

		var res api.JobRegisterResponse

		res.JobModifyIndex = 33

		jsonResp, err := json.Marshal(res)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal 2. Err: %s", err)
		}
		w.Write(jsonResp)
	})

	// Nomad 'get' (of job status) simulation
	mux.HandleFunc("/v1/job/dp-search-data-finder", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Nomad mock: got get job status\n")

		time.Sleep(100 * time.Millisecond) // small delay to simulate work being done

		w.Header().Set("Content-Type", "application/json")

		JobTypeService := "service"
		var res = api.Job{Type: &JobTypeService}

		jsonResp, err := json.Marshal(res)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal 3. Err: %s", err)
		}
		w.Write(jsonResp)
	})

	// Nomad 'get' (of deployment status) simulation
	mux.HandleFunc("/v1/job/dp-search-data-finder/deployments", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Nomad mock: got get deployment status\n")

		time.Sleep(100 * time.Millisecond) // small delay to simulate work being done

		w.Header().Set("Content-Type", "application/json")

		// !!! TODO / try: Only return 'DeploymentStatusSuccessful' once every 3rd time this
		//     is called (sus out what else can be returned on the 2 other occasions),
		//     to better simulate a real system.
		//     To achieve this (when many jobs are running in parallel), i'd have to keep
		//     a list of job id's that are active and their progress & search / update the list.
		//     This extra job progress would need handling for all the Nomad endpoint
		//     simulations in this file.
		var deployments = []api.Deployment{
			{JobSpecModifyIndex: 33, Status: structs.DeploymentStatusSuccessful},
		}

		jsonResp, err := json.Marshal(deployments)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal 4. Err: %s", err)
		}
		w.Write(jsonResp)
	})

	// Vault is handled via "go-chi" 'NewRouter'
	rVault := chi.NewRouter()
	
	// Vault 'healthcheck' simulation
	rVault.Get("/v1/sys/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Vault mock: health")

		w.Header().Set("Content-Type", "application/json")

		// this type copied from vault source code:
		type HealthResponse struct {
			Initialized                bool   `json:"initialized"`
			Sealed                     bool   `json:"sealed"`
			Standby                    bool   `json:"standby"`
			PerformanceStandby         bool   `json:"performance_standby"`
			ReplicationPerformanceMode string `json:"replication_performance_mode"`
			ReplicationDRMode          string `json:"replication_dr_mode"`
			ServerTimeUTC              int64  `json:"server_time_utc"`
			Version                    string `json:"version"`
			ClusterName                string `json:"cluster_name,omitempty"`
			ClusterID                  string `json:"cluster_id,omitempty"`
			LastWAL                    uint64 `json:"last_wal,omitempty"`
		}

		// send a proper Mock health response for checks to be happy
		vaultHealthResponse := HealthResponse{
			Initialized:                true,
			Sealed:                     false,
			Standby:                    false,
			PerformanceStandby:         false,
			ReplicationPerformanceMode: "disabled",
			ReplicationDRMode:          "disabled",
			ServerTimeUTC:              1516639589,
			Version:                    "0.9.2",
			ClusterName:                "vault-cluster-3bd69ca2",
			ClusterID:                  "00af5aa8-c87d-b5fc-e82e-97cd8dfaf731",
		}

		jsonResp, err := json.Marshal(vaultHealthResponse)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal 5. Err: %s", err)
		}
		w.Write(jsonResp)
	})

	// Vault 'secrets write' simulation
	/*
	   Example of incoming 'secret(s)' packet to process:
	   "PUT /v1/secret/babbage-publishing HTTP/1.1\r\nHost: 127.0.0.1:8200\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 590\r\nX-Vault-Request: true\r\nAccept-Encoding: gzip\r\n\r\n{\"CONTENT_SERVICE_URL\":\"http://$NOMAD_IP_http:10800/v1\",\"DP_LOGGING_FORMAT\":\"json\",\"ELASTIC_SEARCH_CLUSTER\":\"cluster\",\"ELASTIC_SEARCH_SERVER\":\"$NOMAD_IP_http\",\"ENABLE_CENSUS_BANNER\":true,\"ENABLE_COVID19_FEATURE\":true,\"ENABLE_SEARCH_SERVICE\":false,\"EXTERNAL_SEARCH_HOST\":\"$NOMAD_IP_http\",\"EXTERNAL_SEARCH_PORT\":11150,\"EXTERNAL_SPELLCHECK_ENABLED\""
	*/
	rVault.Put("/v1/secret/{secretName}", func(w http.ResponseWriter, r *http.Request) {
		secretName := chi.URLParam(r, "secretName")
		log.Println("Vault mock: secret -", secretName)

		time.Sleep(5 * time.Millisecond) // small delay to simulate work being done

		var configSecrets map[string]interface{}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &configSecrets); err != nil {
			log.Printf("Error unbarshaling body: %v", err)
			http.Error(w, "can't read body", http.StatusUnprocessableEntity)
			return
		}
		for k, v := range configSecrets {
			if v2, ok := v.(string); ok {
				if strings.Contains(v2, "BEGIN CERTIFICATE") ||
					strings.Contains(v2, "BEGIN PGP PUBLIC KEY BLOCK") ||
					strings.Contains(v2, "BEGIN PRIVATE KEY") {
					// skip showing sensitive stuff
					continue
				}
			}
			fmt.Printf("key[%s] value[%s]\n", k, v)
			break // only need to show first value to demonstrate that we have the config secrets for an app.
		}

		w.Header().Set("Content-Type", "application/json")

		w.Header().Set("Server", "A Vault Simulation Server")
		w.WriteHeader(200)
	})

	go func() {
		log.Println("Starting Vault Simulation server, utilising go-chi")
		log.Fatal(http.ListenAndServe(":8200", rVault))
	}()

	log.Println("Starting Nomad Simulation server, utilising standard NewServeMux")
	log.Fatal(http.ListenAndServeTLS(":8080", CertPath, KeyPath, mux))
}

//var planResponseMessage = `{ "Index": 0, "NextPeriodicLaunch": "0001-01-01T00:00:00Z", "Warnings": "", "Diff": { "Type": "Added", "TaskGroups": [ { "Updates": { "create": 1 }, "Type": "Added", "Tasks": [ { "Type": "Added", "Objects": ["..."], "Name": "redis", "Fields": [ { "Type": "Added", "Old": "", "New": "docker", "Name": "Driver", "Annotations": null }, { "Type": "Added", "Old": "", "New": "5000000000", "Name": "KillTimeout", "Annotations": null } ], "Annotations": ["forces create"] } ], "Objects": ["..."], "Name": "cache", "Fields": ["..."] } ], "Objects": [ { "Type": "Added", "Objects": null, "Name": "Datacenters", "Fields": ["..."] }, { "Type": "Added", "Objects": null, "Name": "Constraint", "Fields": ["..."] }, { "Type": "Added", "Objects": null, "Name": "Update", "Fields": ["..."] } ], "ID": "example", "Fields": ["..."] }, "CreatedEvals": [ { "ModifyIndex": 0, "CreateIndex": 0, "SnapshotIndex": 0, "AnnotatePlan": false, "EscapedComputedClass": false, "NodeModifyIndex": 0, "NodeID": "", "JobModifyIndex": 0, "JobID": "example", "TriggeredBy": "job-register", "Type": "batch", "Priority": 50, "ID": "312e6a6d-8d01-0daf-9105-14919a66dba3", "Status": "blocked", "StatusDescription": "created to place remaining allocations", "Wait": 0, "NextEval": "", "PreviousEval": "80318ae4-7eda-e570-e59d-bc11df134817", "BlockedEval": "", "FailedTGAllocs": null, "ClassEligibility": { "v1:7968290453076422024": true } } ], "JobModifyIndex": 0, "FailedTGAllocs": { "cache": { "CoalescedFailures": 3, "AllocationTime": 46415, "Scores": null, "NodesEvaluated": 1, "NodesFiltered": 0, "NodesAvailable": { "dc1": 1 }, "ClassFiltered": null, "ConstraintFiltered": null, "NodesExhausted": 1, "ClassExhausted": null, "DimensionExhausted": { "cpu": 1 } } }, "Annotations": { "DesiredTGUpdates": { "cache": { "DestructiveUpdate": 0, "InPlaceUpdate": 0, "Stop": 0, "Migrate": 0, "Place": 11, "Ignore": 0 } } } }`

/*var planResponseMessage = `
{
	"Index": 0,
	"NextPeriodicLaunch": "0001-01-01T00:00:00Z",
	"Warnings": "",
	"Diff": {
	  "Type": "Added",
	  "TaskGroups": [
		{
		  "Updates": {
			"create": 1
		  },
		  "Type": "Added",
		  "Tasks": [
			{
			  "Type": "Added",
			  "Objects": ["..."],
			  "Name": "redis",
			  "Fields": [
				{
				  "Type": "Added",
				  "Old": "",
				  "New": "docker",
				  "Name": "Driver",
				  "Annotations": null
				},
				{
				  "Type": "Added",
				  "Old": "",
				  "New": "5000000000",
				  "Name": "KillTimeout",
				  "Annotations": null
				}
			  ],
			  "Annotations": ["forces create"]
			}
		  ],
		  "Objects": ["..."],
		  "Name": "cache",
		  "Fields": ["..."]
		}
	  ],
	  "Objects": [
		{
		  "Type": "Added",
		  "Objects": null,
		  "Name": "Datacenters",
		  "Fields": ["..."]
		},
		{
		  "Type": "Added",
		  "Objects": null,
		  "Name": "Constraint",
		  "Fields": ["..."]
		},
		{
		  "Type": "Added",
		  "Objects": null,
		  "Name": "Update",
		  "Fields": ["..."]
		}
	  ],
	  "ID": "example",
	  "Fields": ["..."]
	},
	"CreatedEvals": [
	  {
		"ModifyIndex": 0,
		"CreateIndex": 0,
		"SnapshotIndex": 0,
		"AnnotatePlan": false,
		"EscapedComputedClass": false,
		"NodeModifyIndex": 0,
		"NodeID": "",
		"JobModifyIndex": 0,
		"JobID": "example",
		"TriggeredBy": "job-register",
		"Type": "batch",
		"Priority": 50,
		"ID": "312e6a6d-8d01-0daf-9105-14919a66dba3",
		"Status": "blocked",
		"StatusDescription": "created to place remaining allocations",
		"Wait": 0,
		"NextEval": "",
		"PreviousEval": "80318ae4-7eda-e570-e59d-bc11df134817",
		"BlockedEval": "",
		"FailedTGAllocs": null,
		"ClassEligibility": {
		  "v1:7968290453076422024": true
		}
	  }
	],
	"JobModifyIndex": 0,
	"FailedTGAllocs": {
	  "cache": {
		"CoalescedFailures": 3,
		"AllocationTime": 46415,
		"Scores": null,
		"NodesEvaluated": 1,
		"NodesFiltered": 0,
		"NodesAvailable": {
		  "dc1": 1
		},
		"ClassFiltered": null,
		"ConstraintFiltered": null,
		"NodesExhausted": 1,
		"ClassExhausted": null,
		"DimensionExhausted": {
		  "cpu": 1
		}
	  }
	},
	"Annotations": {
	  "DesiredTGUpdates": {
		"cache": {
		  "DestructiveUpdate": 0,
		  "InPlaceUpdate": 0,
		  "Stop": 0,
		  "Migrate": 0,
		  "Place": 11,
		  "Ignore": 0
		}
	  }
	}
  }
`
*/
