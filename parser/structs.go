package parser

type LoadAverage struct {
	Load1  float32
	Load5  float32
	Load15 float32
}

type CpuLoad struct {
	UserLoad   float32
	SystemLoad float32
	IDLELoad   float32
}

type DiskLoad struct {
	Names     []string
	TPS       []float32
	WriteLoad []float32
	ReadLoad  []float32
	CPU       CpuLoad
}

type DiskFree struct {
	Names     []string
	MBFree    []float32
	MBUses    []float32
	INodeFree []float32
	INodeUses []float32
}

type NetStat struct {
	Protocols []string
	Address   []string
	PID       []string
	States    map[string]int32
}

type TopTalkers struct {
	ProtocolBytes map[string]int32
	Source        []string
	Destination   []string
	Protocol      []string
	BPS           []float32
}

type ParsedData struct {
	LoadAverage LoadAverage
	CpuLoad     CpuLoad
	DiskLoad    DiskLoad
	DiskFree    DiskFree
	NetStat     NetStat
	TopTalkers  TopTalkers
}
