package mqtt

import (
	"encoding/json"
	"go_node_engine/model"
	"testing"
)

func TestReportServiceLoadMetrics(t *testing.T) {
	pubs := []struct{ topic, payload string }{}
	old := publishToBrokerFn
	publishToBrokerFn = func(topic, payload string) { pubs = append(pubs, struct{ topic, payload string }{topic, payload}) }
	defer func() { publishToBrokerFn = old }()

	services := []model.Resources{
		{Sname: "svc.a", Instance: 0, Cpu: "0.5", Memory: "0.25"},
		{Sname: "svc.b", Instance: 1, Cpu: "50%", Memory: "75%"},
	}
	ReportServiceLoadMetrics(services)
	found := false
	for _, p := range pubs {
		if p.topic == "jobs/load_metrics" {
			found = true
			t.Logf("Published payload: %s", p.payload)
			var dec struct {
				LoadMetrics []struct {
					JobName           string  `json:"job_name"`
					InstanceNumber    int     `json:"instance_number"`
					CpuUsage          float64 `json:"cpu_usage"`
					MemoryUsage       float64 `json:"memory_usage"`
					ActiveConnections int     `json:"active_connections"`
				} `json:"load_metrics"`
			}
			if err := json.Unmarshal([]byte(p.payload), &dec); err != nil {
				t.Fatalf("json err: %v", err)
			}
			if len(dec.LoadMetrics) != 2 {
				t.Fatalf("expected 2 metrics got %d", len(dec.LoadMetrics))
			}
			t.Logf("Decoded metrics[0]: job=%s inst=%d cpu=%.2f mem=%.2f", dec.LoadMetrics[0].JobName, dec.LoadMetrics[0].InstanceNumber, dec.LoadMetrics[0].CpuUsage, dec.LoadMetrics[0].MemoryUsage)
			t.Logf("Decoded metrics[1]: job=%s inst=%d cpu=%.2f mem=%.2f", dec.LoadMetrics[1].JobName, dec.LoadMetrics[1].InstanceNumber, dec.LoadMetrics[1].CpuUsage, dec.LoadMetrics[1].MemoryUsage)
			if dec.LoadMetrics[0].CpuUsage != 0.5 || dec.LoadMetrics[0].MemoryUsage != 0.25 {
				t.Fatalf("first metric mismatch %+v", dec.LoadMetrics[0])
			}
			if dec.LoadMetrics[1].CpuUsage != 0.5 || dec.LoadMetrics[1].MemoryUsage != 0.75 {
				t.Fatalf("percent parse mismatch %+v", dec.LoadMetrics[1])
			}
		}
	}
	if !found {
		t.Fatalf("jobs/load_metrics not published")
	}
}
