package main

import (
	"flag"
	"github.com/q191201771/lal/pkg/aac"
	"github.com/q191201771/lal/pkg/avc"
	"github.com/q191201771/lal/pkg/httpflv"
	"github.com/q191201771/nezha/pkg/log"
	"io"
	"os"
)

func main() {
	var err error
	flvFileName, aacFileName, avcFileName := parseFlag()

	var ffr httpflv.FlvFileReader
	err = ffr.Open(flvFileName)
	log.FatalIfErrorNotNil(err)
	defer ffr.Dispose()
	log.Infof("open flv file succ.")

	afp, err := os.Create(aacFileName)
	log.FatalIfErrorNotNil(err)
	defer afp.Close()
	log.Infof("open es aac file succ.")

	vfp, err := os.Create(avcFileName)
	log.FatalIfErrorNotNil(err)
	defer vfp.Close()
	log.Infof("open es h264 file succ.")

	_, err = ffr.ReadFlvHeader()
	log.FatalIfErrorNotNil(err)

	for {
		tag, err := ffr.ReadTag()
		if err == io.EOF {
			log.Infof("EOF.")
			break
		}
		log.FatalIfErrorNotNil(err)

		payload := tag.Payload()

		switch tag.Header.T {
		case httpflv.TagTypeAudio:
			aac.CaptureAAC(afp, payload)
		case httpflv.TagTypeVideo:
			avc.CaptureAVC(vfp, payload)
		}
	}
}

func parseFlag() (string, string, string) {
	flv := flag.String("i", "", "specify flv file")
	aac := flag.String("a", "", "specify es aac file")
	avc := flag.String("v", "", "specify es h264 file")
	flag.Parse()
	if *flv == "" || *avc == "" || *aac == "" {
		flag.Usage()
		os.Exit(1)
	}
	return *flv, *aac, *avc
}