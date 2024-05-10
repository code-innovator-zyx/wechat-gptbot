package utils

import (
	bytes2 "bytes"
	"testing"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/4/10 16:17
* @Package:
 */

func TestDownloadImage(t *testing.T) {

	t.Run("original", func(t *testing.T) {
		reader := bytes2.Buffer{}
		err := DownloadImage("https://dalleprodsec.blob.core.windows.net/private/images/fc8fba77-4343-4ae4-97a1-b66e771995e4/generated_00.png?se=2024-04-11T09%3A05%3A11Z&sig=Y%2FG3w8wDWdv%2FHV0NTDOCMp1P5ZCMVDCXJkPgA%2BSM%2BTQ%3D&ske=2024-04-16T12%3A53%3A04Z&skoid=e52d5ed7-0657-4f62-bc12-7e5dbb260a96&sks=b&skt=2024-04-09T12%3A53%3A04Z&sktid=33e01921-4d64-4f8c-a055-5bdaffd5e33d&skv=2020-10-02&sp=r&spr=https&sr=b&sv=2020-10-02", &reader)
		if nil != err {
			t.Error(err)
			return
		} //3163177
		t.Log(len(reader.Bytes())) //149262
	})

	t.Run("compress", func(t *testing.T) {
		reader := bytes2.Buffer{}
		err := CompressImage("https://dalleprodsec.blob.core.windows.net/private/images/fc8fba77-4343-4ae4-97a1-b66e771995e4/generated_00.png?se=2024-04-11T09%3A05%3A11Z&sig=Y%2FG3w8wDWdv%2FHV0NTDOCMp1P5ZCMVDCXJkPgA%2BSM%2BTQ%3D&ske=2024-04-16T12%3A53%3A04Z&skoid=e52d5ed7-0657-4f62-bc12-7e5dbb260a96&sks=b&skt=2024-04-09T12%3A53%3A04Z&sktid=33e01921-4d64-4f8c-a055-5bdaffd5e33d&skv=2020-10-02&sp=r&spr=https&sr=b&sv=2020-10-02", &reader)
		if nil != err {
			t.Error(err)
			return
		}
		t.Log(len(reader.Bytes())) //

	})
	t.Run("image type", func(t *testing.T) {
		GetImageSuffix("https://dalleprodsec.blob.core.windows.net/private/images/fc8fba77-4343-4ae4-97a1-b66e771995e4/generated_00.png?se=2024-04-11T09%3A05%3A11Z&sig=Y%2FG3w8wDWdv%2FHV0NTDOCMp1P5ZCMVDCXJkPgA%2BSM%2BTQ%3D&ske=2024-04-16T12%3A53%3A04Z&skoid=e52d5ed7-0657-4f62-bc12-7e5dbb260a96&sks=b&skt=2024-04-09T12%3A53%3A04Z&sktid=33e01921-4d64-4f8c-a055-5bdaffd5e33d&skv=2020-10-02&sp=r&spr=https&sr=b&sv=2020-10-02")
	})

	t.Run("test", func(t *testing.T) {

	})
}
