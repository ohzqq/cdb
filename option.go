package cdb

type Opt func(*Lib)

func IsAudiobooks() Opt {
	return func(l *Lib) {
		l.isAudiobooks = true
	}
}

//func PrintQuery() Opt {
//  return func(db *Lib) {
//    db.printQuery = true
//  }
//}
