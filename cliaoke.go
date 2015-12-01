package main

import (
	"fmt"
	"os"

	"os/signal"
)

func usage() {
	fmt.Println("cli-aoke")
	fmt.Println("Command Line Karaoke")
	fmt.Println("Version 1.1.0")
	fmt.Println("Usage ===>")
	fmt.Println("           songs: lists the song files you can choose from")
	fmt.Println("           sing:  starts a song!")
}

func main() {
	if len(os.Args) <= 1 {
		usage()
		return
	}
	switch os.Args[1] {
	case "songs":
		err := songs_get()
		if err != nil {
			fmt.Println(err)
		}
	case "sing":
		if len(os.Args) <= 2 {
			fmt.Println("You did not pass a song file name.")
			fmt.Println("Run `cli-aoke songs` to get the list.")
			fmt.Println("Then run `cli-aoke sing <song_file_name>`.")
			return
		}

		song_file := os.Args[2]
		if !song_exists(song_file) {
			fmt.Println("Sowwwwyyyyy")
			fmt.Println(song_file, "is not in the ~/.cli-aoke directory")
			return
		}

		fmt.Println("DJ cli-aoke on the request line.")
		fmt.Println("Bringing up the track...")

		fmt.Println("Fetching lyrics...")
		song_lyrics := lyrics_get(song_file)
		if song_lyrics == nil {
			fmt.Println("Sowwwwyyyyy")
			fmt.Println("Search for", song_file, "returned zero results.")
			return
		}

		quitchan := make(chan os.Signal, 1)
		signal.Notify(quitchan, os.Interrupt)
		// kar := karaoke(song_file, song_lyrics)
	default:
		usage()
	}
}
