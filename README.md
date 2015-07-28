 

This is a simple wrapper with REST interface for [SoX](http://sox.sourceforge.net/), particularly it's spectrogram tool, written in Golang
## Usage
Make a POST request…
```bash
$ curl -F "file=@some-audio.mp3" -F "file=@another-audio.wav" host:port/api/v1/post
```
and you will get a JSON response:
```JSON
{"results":
	{
		"some-audio.mp3":"host:port/png/some-audio.mp3.png",
		"another-audio.wav":"host:port/png/some-audio.wav.png"
	}
}


```
**Note:** This was written mainly for educational purposes, if you have any comments or suggestions, feel free to drop me a line to <denis.kolosov@gmail.com>

## Dependencies
You will need SoX to get actual spectrograms. Install on Debian-like systems:
```bash
$  sudo apt-get install sox
```
If you need to convert mp3s:
```bash
$  sudo apt-get install sox libsox-fmt-mp3 
```


## API Reference
Currently (and probably it won't change) there is only one endpoint: /api/v1/post

## Tests
TODO: Write some tests
