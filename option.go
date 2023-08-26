package calibredb

type Opt func(*DB)

func IsAudiobooks() Opt {
	return func(db *DB) {
		db.isAudiobooks = true
	}
}

func PrintQuery() Opt {
	return func(db *DB) {
		db.printQuery = true
	}
}
