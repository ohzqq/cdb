package cdb

type Option func(*Lib)

func IsAudiobooks() Option {
	return func(l *Lib) {
		l.isAudiobooks = true
	}
}

//func PrintQuery() Opt {
//  return func(db *Lib) {
//    db.printQuery = true
//  }
//}
