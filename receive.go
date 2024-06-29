package wa
import (
    //"encoding/json"
    //"fmt"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
    "strings"
    waProto "go.mau.fi/whatsmeow/binary/proto"
)

type M struct {
    From types.JID
    Sender types.JID
    ID string
    PushName string
    FromMe bool
    IsGroup bool
    IsOwner bool
    Text string
    Args []string
    Command string
    Prefix string
    Q string
    HasQuoted bool
    Quoted struct {
        Msg *waProto.Message
        Image *waProto.ImageMessage
        Video *waProto.VideoMessage
        Sticker *waProto.StickerMessage
    }
    Full *events.Message
}

func Terima(msg *events.Message) *M {
    var m = &M{}

    m.Full = msg
    m.ID = msg.Info.ID
    m.From = msg.Info.Chat
    m.FromMe = msg.Info.IsFromMe
    m.IsGroup = msg.Info.IsGroup
    m.Sender = msg.Info.Sender
    m.PushName = msg.Info.PushName

    extended := msg.Message.GetExtendedTextMessage()
    extendedT := extended.GetText()
    text := msg.Message.GetConversation()
    imageMatch := msg.Message.GetImageMessage().GetCaption()
    videoMatch := msg.Message.GetVideoMessage().GetCaption()
    var body string
    if text != "" {
        body = text
    } else if imageMatch != "" {
        body = imageMatch
    } else if videoMatch != "" {
        body = videoMatch
    } else if extendedT != "" {
        body = extendedT
    }
    m.Text = body
    m.Args = strings.Split(body, " ")
    m.Prefix = "."
    m.Command = strings.ToLower(m.Args[0])
    m.Prefix = "."
    m.Q = strings.Join(m.Args[1:], ` `)
    
    m.Quoted.Msg = extended.GetContextInfo().GetQuotedMessage()
    m.Quoted.Image = m.Quoted.Msg.GetImageMessage()
    m.Quoted.Video = m.Quoted.Msg.GetVideoMessage()
    m.Quoted.Sticker = m.Quoted.Msg.GetStickerMessage()
    return m
}