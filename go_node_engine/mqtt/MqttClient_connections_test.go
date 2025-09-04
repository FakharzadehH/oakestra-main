package mqtt

import (
	"strings"
	"testing"
	"go_node_engine/model"
)

func TestReportServiceLoadMetricsIncludesConnections(t *testing.T) {
	var captured string
	old := publishToBrokerFn
	publishToBrokerFn = func(topic, payload string) { if topic == "jobs/load_metrics" { captured = payload } }
	defer func(){ publishToBrokerFn = old }()

	ReportServiceLoadMetrics([]model.Resources{{Cpu: "10%", Memory: "20%", Sname: "svc.a", Instance: 0}})
	if captured == "" { t.Fatalf("expected payload published") }
	if !strings.Contains(captured, "active_connections") { t.Fatalf("missing active_connections field: %s", captured) }
}

func TestCountActiveTcpConnections(t *testing.T) {
	c := countActiveTcpConnections()
	if c < -1 { t.Fatalf("unexpected value %d", c) }
}
