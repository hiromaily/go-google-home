package ttsengine

//import (
//	"fmt"
//	"os"
//
//	tts "github.com/pqyptixa/tts2media"
//)
//
//var dataPath = os.Getenv("GOPATH") + "/src/github.com/hiromaily/go-google-home/data/"
//
////type EspeakSpeech struct {
////	Text     string
////	Lang     string
////	Speed    string
////	Gender   string
////	Altvoice string
////	Quality  string
////	Pitch    string
////}
//
//func createWAV(text string) (*tts.Media, error) {
//	espwav := &tts.EspeakSpeech{Text: text, Lang: "en", Speed: "135", Gender: "m", Altvoice: "0", Quality: "high", Pitch: "50"}
//
//	media, err := espwav.NewEspeakSpeech()
//	if err != nil {
//		return nil, fmt.Errorf("NewEspeakSpeech returned error: %v", err)
//	}
//
//	filename := dataPath + media.Filename + ".wav"
//
//	if _, err = os.Stat(filename); err != nil {
//		return nil, fmt.Errorf("There was an error opening the WAV file: %v", err)
//	}
//	return media, nil
//}
//
//func removeWAV(media *tts.Media) {
//	media.RemoveWAV()
//}
