package bpm

type InterfaceBpm interface {
	ExtractBpm(fileUrl string) (float64, bool, error)
}

//func Bpm(
//	urls []string,
//	audio InterfaceBpm,
//	//destination []datastore.TrackUpdateRequest,
//	basePath string,
//	) (*modules.TrackBpmResult, error) {
//
//	for _, v := range urls {
//		audio.ExtractBpm(v)
//	}
//
//	return nil
//}


