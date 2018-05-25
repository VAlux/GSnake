package highscore

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	t "time"
)

const highScoreFilename = "score.hsc"

// HighScore represents all of the single high-score entry components
type HighScore struct {
	Timestamp  t.Time
	Score      int
	PlayerName string
}

func initialize() {
	gob.Register(HighScore{})
}

func serialize(score *HighScore) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(*score)
	if err != nil {
		log.Panic("Error serialization high score:", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

func deSerialize(file *os.File) (*HighScore, error) {
	score := &HighScore{}
	decoder := gob.NewDecoder(file)
	err := decoder.Decode(score)
	if err != nil {
		log.Panic("Error de-serializing high score:", err)
		return nil, err
	}
	return score, nil
}

// Save writes high score structure to file
func Save(score *HighScore) {
	payload, serializeError := serialize(score)
	if serializeError != nil {
		log.Panic("Error saving high score to file:", serializeError)
		return
	}

	saveError := ioutil.WriteFile(highScoreFilename, payload, 0666)
	if saveError != nil {
		log.Panic("Error saving high score to file:", saveError)
		return
	}

	log.Printf("High score successfully saved to file: %s", highScoreFilename)
}

// Load reads high score structure from file
func Load() (*HighScore, error) {
	log.Printf("Loading high score from file: %s", highScoreFilename)
	file, readError := os.Open(highScoreFilename)
	if readError != nil {
		return nil, readError
	}

	score, deSerializeError := deSerialize(file)
	if deSerializeError != nil {
		return nil, readError
	}

	return score, nil
}
