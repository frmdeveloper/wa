package wa
import (
    "encoding/base64"
    "fmt"
    "github.com/google/uuid"
    "io/ioutil"
    "net/http"
    "os"
)

type Getbyte struct {
    Byte []byte
    Mimetype string
    Length int
}

func Atob(base string) ([]byte) {
    /*var decoded = make([]byte, base64.StdEncoding.DecodedLen(len(args.Base64)))
        var _, err = base64.StdEncoding.Decode(decoded, []byte(args.Base64))
        if err != nil { return nil, err }
        return decoded,nil*/
    b,_ := base64.StdEncoding.DecodeString(base)
    return b
}
func Btoa(buffer []byte) string {
  return base64.StdEncoding.EncodeToString(buffer)
}
func GetByte(args L) (*Getbyte, error) {
    if args.Byte != nil {
        return &Getbyte{
            Byte: args.Byte, 
            Mimetype: http.DetectContentType(args.Byte),
            Length: len(args.Byte),
        }, nil
    }
    if args.Filepath != "" {
        bacaf,erbacaf := os.ReadFile(args.Filepath)
        return &Getbyte{
            Byte: bacaf, 
            Mimetype: http.DetectContentType(bacaf),
            Length: len(bacaf),
        }, erbacaf
    }
    if args.Url != "" {
        res, err := http.Get(args.Url)
        if err != nil { return nil, err }
        defer res.Body.Close()
        rio,errio := ioutil.ReadAll(res.Body)
        return &Getbyte{
            Byte: rio,
            Mimetype: http.DetectContentType(rio),
            Length: len(rio),
        }, errio
    }
    if args.Base64 != "" {
        rtob := Atob(args.Base64)
        if rtob == nil {
            return nil,fmt.Errorf("error base64")
        }
        return &Getbyte{
            Byte: rtob,
            Mimetype: http.DetectContentType(rtob),
            Length: len(rtob),
        }, nil
    }
    if args.Text != "" {
        tobyte := []byte(args.String)
        return &Getbyte{
            Byte: tobyte, 
            Mimetype: http.DetectContentType(tobyte),
            Length: len(tobyte),
        }, nil
    }
    return nil,nil
}
func U() string {
  return fmt.Sprintf("%s",uuid.New())
}