package bpm

type InterfaceBpm interface {
	ExtractBpm(fileUrl string) (float64, bool, error)
}


