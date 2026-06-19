package models

// Metrics is the full catalog of telemetry fields. Entries with an empty
// Types slice are universal (every machine reports them). Entries with Types
// set are only returned for machines whose type is listed.
var Metrics = []MetricDef{
	// Universal
	{Key: "output", Label: "Output Count", Unit: "pcs", Min: 0, Max: 10000},
	{Key: "status", Label: "Machine Status", Unit: "", Min: 0, Max: 0},

	// Shared across several types
	{Key: "temperature", Label: "Temperature", Unit: "°C", Min: 20, Max: 95,
		Types: []string{"CNC", "Injection Molding", "Extruder", "Press", "Mixer", "Conveyor", "Robot", "Filler"}},
	{Key: "pressure", Label: "Pressure", Unit: "bar", Min: 0, Max: 10,
		Types: []string{"Injection Molding", "Press", "Filler", "Extruder"}},
	{Key: "speed", Label: "Speed", Unit: "rpm", Min: 0, Max: 3000,
		Types: []string{"CNC", "Mixer", "Conveyor", "Extruder", "Robot", "Labeler"}},
	{Key: "vibration", Label: "Vibration", Unit: "mm/s", Min: 0, Max: 2,
		Types: []string{"CNC", "Conveyor", "Robot", "Press", "Extruder", "Injection Molding", "Mixer", "Filler", "Labeler"}},

	// CNC
	{Key: "spindle_load", Label: "Spindle Load", Unit: "%", Min: 0, Max: 100,
		Types: []string{"CNC"}},
	{Key: "tool_wear", Label: "Tool Wear", Unit: "%", Min: 0, Max: 100,
		Types: []string{"CNC"}},

	// Conveyor
	{Key: "belt_tension", Label: "Belt Tension", Unit: "N", Min: 50, Max: 500,
		Types: []string{"Conveyor"}},

	// Robot
	{Key: "joint_torque", Label: "Joint Torque", Unit: "Nm", Min: 0, Max: 200,
		Types: []string{"Robot"}},
	{Key: "payload", Label: "Payload", Unit: "kg", Min: 0, Max: 150,
		Types: []string{"Robot"}},

	// Press
	{Key: "force", Label: "Force", Unit: "kN", Min: 0, Max: 500,
		Types: []string{"Press"}},
	{Key: "stroke_count", Label: "Stroke Count", Unit: "strokes", Min: 0, Max: 50000,
		Types: []string{"Press"}},

	// Mixer
	{Key: "torque", Label: "Torque", Unit: "Nm", Min: 0, Max: 300,
		Types: []string{"Mixer"}},

	// Filler
	{Key: "fill_rate", Label: "Fill Rate", Unit: "units/min", Min: 0, Max: 120,
		Types: []string{"Filler"}},

	// Labeler
	{Key: "reject_rate", Label: "Reject Rate", Unit: "%", Min: 0, Max: 15,
		Types: []string{"Labeler"}},
}

// MetricsForType returns the metrics that apply to the given machine type:
// universal metrics (empty Types) plus those whose Types list contains machineType.
// Unknown types get only the universal metrics.
func MetricsForType(machineType string) []MetricDef {
	result := make([]MetricDef, 0, len(Metrics))
	for _, m := range Metrics {
		if len(m.Types) == 0 {
			result = append(result, m)
			continue
		}
		for _, t := range m.Types {
			if t == machineType {
				result = append(result, m)
				break
			}
		}
	}
	return result
}

// MetricByKey looks up a metric definition by its key.
func MetricByKey(key string) (MetricDef, bool) {
	for _, m := range Metrics {
		if m.Key == key {
			return m, true
		}
	}
	return MetricDef{}, false
}
