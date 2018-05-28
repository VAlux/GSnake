package highscore

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	t "time"
)

const highScoreFilename = "score.hsc"
const key = "cegthctrm.hysqrk.xrjnjhsqytdjpvj"

// HighScore represents all of the single high-score entry components
type HighScore struct {
	Timestamp  t.Time
	Score      int
	PlayerName string
}

// HighScores represents a slice of HighScore entries
type HighScores []HighScore

func initialize() {
	gob.Register(HighScore{})
}

func (score *HighScore) String() string {
	return score.Timestamp.String() + " :: " + score.PlayerName + " :: " + string(score.Score)
}

func (scores *HighScores) String() string {
	content := ""
	for _, score := range *scores {
		content += score.String() + "\n"
	}
	return content
}

func serialize(scores *HighScores) ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(*scores)
	if err != nil {
		log.Panic("Error serialization high score:", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

func deSerialize(file *os.File) (*HighScores, error) {
	content, readingError := ioutil.ReadAll(file)
	if readingError != nil {
		log.Panic("Error reading contents from file:", readingError)
		return nil, readingError
	}

	decrypted, decryptionError := decrypt(content)
	if decryptionError != nil {
		log.Panic("Error decrypting high score contents:", decryptionError)
		return nil, decryptionError
	}

	score := new(HighScores)
	decoder := gob.NewDecoder(bytes.NewReader(decrypted))
	err := decoder.Decode(score)
	if err != nil {
		log.Panic("Error de-serializing high score:", err)
		return nil, err
	}
	return score, nil
}

// Save writes high score structure to file
func Save(score *HighScore) {
	currentScores, _ := Load()
	currentScores = append(currentScores, *score)

	log.Println("Saving current high scores: " + currentScores.String())

	payload, serializeError := serialize(&currentScores)
	if serializeError != nil {
		log.Panic("Error saving high score to file:", serializeError)
		return
	}

	encryptedPayload, encryptionError := encrypt(payload)
	if encryptionError != nil {
		log.Panic("Error encryption of the high score payload:", encryptionError)
	}

	saveError := ioutil.WriteFile(highScoreFilename, encryptedPayload, 0666)
	if saveError != nil {
		log.Panic("Error saving high score to file:", saveError)
		return
	}

	log.Printf("High score successfully saved to file: %s", highScoreFilename)
}

// Load reads high score structure from file
func Load() (HighScores, error) {
	log.Printf("Loading high score from file: %s", highScoreFilename)
	file, readError := os.Open(highScoreFilename)
	if readError != nil {
		return nil, readError
	}

	scores, deSerializeError := deSerialize(file)
	if deSerializeError != nil {
		return nil, readError
	}

	return *scores, nil
}

func encrypt(payload []byte) ([]byte, error) {
	keyBytes := []byte(key)
	block, chipherCreationError := aes.NewCipher(keyBytes)
	if chipherCreationError != nil {
		log.Panic("Error creating cipher:", chipherCreationError)
		return nil, chipherCreationError
	}

	encryptedPayload := make([]byte, aes.BlockSize+len(payload))
	iv := encryptedPayload[:aes.BlockSize]

	_, ivGenerationError := io.ReadFull(rand.Reader, iv)
	if ivGenerationError != nil {
		log.Panic("Error generating IV:", ivGenerationError)
		return nil, ivGenerationError
	}

	encryptionStream := cipher.NewCFBEncrypter(block, iv)
	encryptionStream.XORKeyStream(encryptedPayload[aes.BlockSize:], payload)

	return encryptedPayload, nil
}

func decrypt(payload []byte) ([]byte, error) {
	keyBytes := []byte(key)
	block, chipherCreationError := aes.NewCipher(keyBytes)
	if chipherCreationError != nil {
		log.Panic("Error creating cipher:", chipherCreationError)
		return nil, chipherCreationError
	}

	if len(payload) < aes.BlockSize {
		errorMessage := "High score file is too short"
		log.Panic("Decryption error:", errorMessage)
		return nil, errors.New(errorMessage)
	}

	iv := payload[:aes.BlockSize]
	encryptedPayload := payload[aes.BlockSize:]

	decryptionStream := cipher.NewCFBDecrypter(block, iv)
	decryptionStream.XORKeyStream(encryptedPayload, encryptedPayload)

	return encryptedPayload, nil
}
