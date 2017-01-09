package gotefs

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

func Serialize(v interface{}) ([]byte, error) {
	var b = bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(v)
	if err != nil {
		log.Printf("error: could not encode data (%v)", err)
	}
	return b.Bytes(), err
}

func SerializeToFile(f *os.File, v interface{}) error {
	return nil
}

func Deserialize(data []byte) (interface{}, error) {
	return nil, nil
}

func DeserializeFile(f *os.File) (interface{}, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("error: could not open file for deserialization (%v)", err)
		return nil, err
	}
	return Deserialize(b)
}
