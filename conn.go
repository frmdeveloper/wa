package wa

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "google.golang.org/protobuf/proto"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
    "regexp"
    "strings"
    waProto "go.mau.fi/whatsmeow/binary/proto"
)

type Mess struct {
    Edit string
    Mentions []string
    ParseMention bool
    Quoted *events.Message
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
func (c *Conn) RelayMessage(jid string, message *waProto.Message, a Mess) (*events.Message, error) {
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
func (c *Conn) SendText(jid string, text string, a Mess) (*events.Message, error) {
    var mentionedjid []string
    if a.ParseMention {
        mentionedjid = c.ParseMention(text)
    } else {
        mentionedjid = a.Mentions
    }
    contekinfo := &waProto.ContextInfo{
        MentionedJID:  mentionedjid,
        StanzaID:      &a.Quoted.Info.ID,
        Participant:   proto.String(a.Quoted.Info.Sender.String()),
        QuotedMessage: a.Quoted.Message,
    }
    return c.RelayMessage(jid, &waProto.Message{
        ExtendedTextMessage: &waProto.ExtendedTextMessage{
            Text: &text,
            ContextInfo: contekinfo,
        },
    }, a)
}