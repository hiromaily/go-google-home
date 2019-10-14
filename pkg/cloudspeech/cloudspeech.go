package cloudspeech

//// Convert sound data to text...
//// https://cloud.google.com/speech/docs/reference/libraries#client-libraries-install-go
//// https://github.com/googleapis/google-cloud-go/blob/master/speech/apiv1/speech_client_example_test.go
//
//import (
//	"fmt"
//	"io/ioutil"
//
//	speech "cloud.google.com/go/speech/apiv1"
//	"golang.org/x/net/context"
//	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
//
//	lg "github.com/hiromaily/golibs/log"
//)
//
//var audioPath = "~/work/audio.raw"
//
//func cloudSpeech() error {
//	ctx := context.Background()
//
//	// Creates a client.
//	client, err := speech.NewClient(ctx)
//	if err != nil {
//		lg.Errorf("Failed to create client: %v", err)
//		//google: could not find default credentials.
//		//https://developers.google.com/accounts/docs/application-default-credentials
//		//https://cloud.google.com/docs/authentication/production
//		return err
//	}
//
//	// Sets the name of the audio file to transcribe.
//	//audioPath
//
//	// Reads the audio file into memory.
//	data, err := ioutil.ReadFile(audioPath)
//	if err != nil {
//		lg.Errorf("Failed to read file: %v", err)
//		return err
//	}
//
//	// Detects speech in the audio file.
//	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
//		Config: &speechpb.RecognitionConfig{
//			Encoding:        speechpb.RecognitionConfig_LINEAR16,
//			SampleRateHertz: 16000,
//			LanguageCode:    "en-US",
//		},
//		Audio: &speechpb.RecognitionAudio{
//			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
//		},
//	})
//	if err != nil {
//		lg.Errorf("Failed to recognize: %v", err)
//		return err
//	}
//
//	// Prints the results.
//	for _, result := range resp.Results {
//		for _, alt := range result.Alternatives {
//			fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
//		}
//	}
//
//	return nil
//}
