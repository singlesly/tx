package p2p

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// Функция для сжатия данных
func compressData(data []byte) []byte {
	var compressedData bytes.Buffer
	gz := gzip.NewWriter(&compressedData)
	_, err := gz.Write(data)
	if err != nil {
		panic(err)
	}
	gz.Close()
	return compressedData.Bytes()
}

// Функция для распаковки данных
func decompressData(data []byte) []byte {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	decompressedData, err := ioutil.ReadAll(gz)
	if err != nil {
		panic(err)
	}
	return decompressedData
}
