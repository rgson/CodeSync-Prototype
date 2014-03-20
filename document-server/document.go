type Document struct {
	shadow			shadow
	backup			backup
}

type shadow struct {
	content			string
	localVersion	int
	remoteVersion	int
}

type backup struct {
	content			string
	localVersion	int
}

func (d *Document) Initialize() {

	c.shadow = shadow{content: "", localVersion: 0, remoteVersion: 0}
	c.backup = backup{content: "", localVersion: 0}
	
}
