package spotify

func InitAll(configfile string) {
	InitConfig(configfile)
	InitDB()
	InitAuth()
	InitSeenWords()
	InitSeenPlaylists()
}
