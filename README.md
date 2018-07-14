# README #

### myRecognize ###
myRecognize is a RESTful API written in GO that accpect .mp3 file, convert it to .raw using SoX and then call Google STT API and response the output.


### Set up ###

#### Dependencies
* Install SoX - a cross-platform (Windows, Linux, MacOS X, etc.) command line utility that can convert various formats of computer audio files in to other formats.
```sh
//Ubuntu
sudo apt-get install sox

//OSX
brew install sox
```

#### Go

* Download Go and setup Go root. Please see https://golang.org/doc/install#download
* Create a Workspace for you go code.  Please see https://golang.org/doc/code.html#Workspaces.
* Set you GOPATH environment variable for your Go root https://golang.org/doc/code.html#GOPATH.
* Run go get github.com/aphpbonn/myRecognize to pull down the repository
* Copy the servapi.env.example to servapi.env and specify path to the Google credential .json file


#### Running the app:


```sh
go run main.go --env-file servapi.env
```

* the application will start on port 8000
* use postman or the following example cURL
```sh
curl -X POST \
  http://localhost:8000/messages/voice \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
	"filename":"/Users/PATH/TO/INPUTFILE.mp3"
}
'
```
