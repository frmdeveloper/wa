package wa
import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "google.golang.org/protobuf/proto"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
    "net/http"
    "regexp"
    "strings"
    waProto "go.mau.fi/whatsmeow/binary/proto"
)

type L struct {
    Edit string
    Caption string
    Mentions []string
    ParseMention bool
    Quoted *events.Message
    Url string
    Base64 string
    Filepath string
    Byte []byte
    Text string
    String string
}
type Conn struct {
    C *whatsmeow.Client
}
func Sends(Cli *whatsmeow.Client) Conn {
    return Conn{
        C: Cli,
    }
}
func (c *Conn) GenerateMessageID() types.MessageID {
    id := make([]byte, 14)
    _, err := rand.Read(id)
    if err != nil {
        panic(err)
    }
    return strings.ToUpper(hex.EncodeToString(id)) + "-FRM"
}
func (c *Conn) ParseMention(text string) []string {
    res := []string{}
    matches := regexp.MustCompile("@([0-9]{5,16}|0)").FindAllStringSubmatch(text, -1)
    for _, match := range matches {
        res = append(res, match[1]+"@s.whatsapp.net")
    }
    return res
}
func (c *Conn) WaUpload(args L, tipeM whatsmeow.MediaType) (*waProto.ImageMessage, error) {
    dow, err := GetByte(args)
    if err != nil { return nil, err }
    uploaded, err := c.C.Upload(context.Background(), dow.Byte, tipeM)
    if err != nil { return nil, err }
    return &waProto.ImageMessage{
        URL:           proto.String(uploaded.URL),
        DirectPath:    proto.String(uploaded.DirectPath),
        MediaKey:      uploaded.MediaKey,
        Mimetype:      proto.String(http.DetectContentType(dow.Byte)),
        FileEncSHA256: uploaded.FileEncSHA256,
        FileSHA256:    uploaded.FileSHA256,
        FileLength:    proto.Uint64(uint64(dow.Length)),
    }, nil
}
func (c *Conn) RelayMessage(jid string, message *waProto.Message, a L) (*events.Message, error) {
    Jid, _ := types.ParseJID(jid)
    if a.Edit != "" {
        message = c.C.BuildEdit(Jid, a.Edit, message)
    }
    send, err := c.C.SendMessage(context.Background(), Jid, message, whatsmeow.SendRequestExtra{ID:c.C.GenerateMessageID()})
    return &events.Message{
        Info: types.MessageInfo{
            ID: send.ID,
            MessageSource: types.MessageSource{
                Chat: Jid,
                Sender: *c.C.Store.ID,
                IsFromMe: true,
                IsGroup: Jid.Server == types.GroupServer,
            },
        },
        Message: message,
    }, err
}
func (c *Conn) quoted(a L) *waProto.ContextInfo {
    return &waProto.ContextInfo{
        StanzaID:      &a.Quoted.Info.ID,
        Participant:   proto.String(a.Quoted.Info.Sender.String()),
        QuotedMessage: a.Quoted.Message,
    }
}
func (c *Conn) SendText(jid string, text string, a L) (*events.Message, error) {
    var mentionedjid []string
    if a.ParseMention {
        mentionedjid = c.ParseMention(text)
    } else {
        mentionedjid = a.Mentions
    }
    co := c.quoted(a)
    co.MentionedJID = mentionedjid
    return c.RelayMessage(jid, &waProto.Message{
        ExtendedTextMessage: &waProto.ExtendedTextMessage{
            Text: &text,
            ContextInfo: co,
        },
    }, a)
}
func (c *Conn) SendImage(jid string, a L) (*events.Message, error) {
    var mentionedjid []string
    if a.ParseMention {
        mentionedjid = c.ParseMention(a.Caption)
    } else {
        mentionedjid = a.Mentions
    }
    up, err := c.WaUpload(a, whatsmeow.MediaImage)
    if err != nil { return nil, err }
    co := c.quoted(a)
    co.MentionedJID = mentionedjid
    return c.RelayMessage(jid, &waProto.Message{
        ImageMessage: &waProto.ImageMessage{
            URL:           up.URL,
            DirectPath:    up.DirectPath,
            MediaKey:      up.MediaKey,
            Mimetype:      up.Mimetype,
            FileEncSHA256: up.FileEncSHA256,
            FileSHA256:    up.FileSHA256,
            FileLength:    up.FileLength,
            Caption:       &a.Caption,
            ContextInfo: co,
        },
    }, a)
}