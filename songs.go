package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

func song_exists(song_file string) bool {
	user, err := user.Current()
	if err != nil {
		fmt.Println("Error looking up user:", err)
		return false
	}
	st, err := os.Stat(filepath.Join(user.HomeDir, ".cliaoke", "songs", song_file))
	if err != nil {
		return false
	}
	if st.IsDir() {
		return false
	}
	return true
}

func songs_get() error {
	user, err := user.Current()
	if err != nil {
		return errors.New("Error looking up user: " + err.Error())
	}
	cliaoke_dir := filepath.Join(user.HomeDir, ".cliaoke")
	_, err = os.Stat(cliaoke_dir)
	if err != nil {
		err = os.MkdirAll(filepath.Join(cliaoke_dir, "songs"), 0700)
		if err != nil {
			return errors.New("Couldn't create songs directory: " + err.Error())
		}
		err = os.MkdirAll(filepath.Join(cliaoke_dir, "lyrics"), 0700)
		if err != nil {
			return errors.New("Couldn't create lyrics directory: " + err.Error())
		}
		err = scrape(filepath.Join(cliaoke_dir, "songs"))
		if err != nil {
			return err
		}
	}

	f, err := os.Open(filepath.Join(cliaoke_dir, "songs"))
	if err != nil {
		return errors.New("Couldn't open songs directory: " + err.Error())
	}
	files, err := f.Readdir(-1)
	for _, fi := range files {
		fmt.Println(fi.Name())
	}
	return err
}

func songs_play(song_file string) error {
	sound_font_file := "/usr/local/share/fluidsynth/generaluser.v.1.44.sf2"
	user, err := user.Current()
	if err != nil {
		return errors.New("Error looking up user: " + err.Error())
	}
	song_path := filepath.Join(user.HomeDir, ".cliaoke", song_file)

	st, err := os.Stat(sound_font_file)
	if err != nil || st.IsDir() {
		return errors.New("You have not installed fluidsynth correctly.\nPlease refer to https://github.com/jfrazelle/cli-aoke.")
	} else if st, err = os.Stat(song_path); err != nil || st.IsDir() {
		return errors.New(song_file + " does not exist")
	}

	proc := exec.Command("fluidsynth", "-i", sound_font_file, song_path)
	proc.Stderr = os.Stdout
	return proc.Run()
}
