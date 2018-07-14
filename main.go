package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"io/ioutil"

	"golang.org/x/net/context"

	"cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"fmt"
	"encoding/json"
	"flag"
	"github.com/joho/godotenv"

	"os/exec"
	"os"
)

var flags struct {
	envFile string
}

// init is the entry point for the entire web application.
func init() {
	log.Println("Starting wwwgo ...")
	flag.StringVar(&flags.envFile, "env-file", "", "Source environment variable from a file")
	flag.Parse()
	if flags.envFile != "" {
		log.Println("Use .env file")
		if err := godotenv.Load(flags.envFile); err != nil {
			log.Fatal("Unable to load .env file: " + err.Error())
		}
	}
}

//main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/messages/voice", RecognizeHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}

type RecognizeRequest struct {
	Filename string `json:"filename"`
}

type RecognizeResponse struct {
	Text string `json:"text"`
}

//RecognizeHandler accepts mp3 file and out put a response
func RecognizeHandler(w http.ResponseWriter, r *http.Request) {

	var req RecognizeRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error Reading request body", err)
		return
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("Error unmarshaling the request", err)
		return
	}


	res, err := recognize(req.Filename)
	if err != nil {
		log.Println(err)
	}

	b, err := json.Marshal(res)
	if err != nil {
		log.Println("Error marshaling the response", err)
		return
	}

	output := string(b)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, output)
}

func convertMP3toLINEAR16(file string) (string, error) {
	out := os.Getenv("GOPATH")+"/src/github.com/aphpbonn/myRecognize/audio.raw"
	cmdArguments := []string{file, "--channels=1", "--rate", "16k", "--bits", "16",
		out}

	cmd := exec.Command("sox", cmdArguments...)
	err := cmd.Run()
	if err != nil {
		log.Fatal("Error converting .mp3 file to .raw file using sox")
		return "",err
	}
	return out,nil
}

func recognize(file string) (RecognizeResponse, error) {

	//convert input file which is .mp3 to .raw format
	converted, err := convertMP3toLINEAR16(file)

	//Prepare calling google cloud speech API
	res := RecognizeResponse{}
	ctx := context.Background()

	// Creates a client.
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return res, err
	}

	// Reads the audio file into memory.
	data, err := ioutil.ReadFile(converted)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
		return res, err
	}

	// Detects speech in the audio file.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_LINEAR16,
			SampleRateHertz: 16000,
			LanguageCode:    "en-US",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})

	if err != nil {
		log.Fatalf("failed to recognize: %v", err)
		return res, err
	}
	// Collect the result
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			res.Text += fmt.Sprintf("%v ", alt.Transcript)
		}
	}
	return res, err
}
