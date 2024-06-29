package wa
import (
    "fmt"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store/sqlstore"
    "go.mau.fi/whatsmeow/types"
    "go.mau.fi/whatsmeow/types/events"
    waLog "go.mau.fi/whatsmeow/util/log"
    //_ "github.com/mattn/go-sqlite3"
    _ "github.com/ncruces/go-sqlite3/driver"
    _ "github.com/ncruces/go-sqlite3/embed"
    //_ "github.com/lib/pq"
)

type Ev struct {
    ChatPresence *events.ChatPresence
    Message *events.Message
    Receipt *events.Receipt
    More interface {}
}
var Conns = make(map[string]*whatsmeow.Client)
func Connect(nomor string, cb func(conn *whatsmeow.Client, evt Ev)) {
    dbLog := waLog.Stdout("Database", "ERROR", true)
    container, err := sqlstore.New("sqlite3", "file:"+nomor+".db?_foreign_keys=on", dbLog)
    //container, err := sqlstore.New("postgres", "postgresql://litegix:qSlCuAVnQXaY@frm1-dfc8ce3933.onlitegix.com:31744?sslmode=disable", dbLog)
    if err != nil { fmt.Println("GoError:",err); return }
    deviceStore, err := container.GetFirstDevice()
    if err != nil { fmt.Println("GoError:",err); return }
    clientLog := waLog.Stdout("Client", "ERROR", true)
    client := whatsmeow.NewClient(deviceStore, clientLog)
    Conns[nomor] = client
    client.AddEventHandler(func(evt interface {}) {
      switch evt.(type) {
  	  case *events.Message:
	      cb(client, Ev{Message:evt.(*events.Message)})
	    case *events.Receipt:
	      cb(client, Ev{Receipt:evt.(*events.Receipt)})
	    case *events.ChatPresence:
	      cb(client, Ev{ChatPresence:evt.(*events.ChatPresence)})
        default:
          cb(client, Ev{More:evt})
      }
    })
    if client.Store.ID == nil {
        err = client.Connect()
        if err != nil { fmt.Println("GoError:",err); return }
        linkingCode, gagal := client.PairPhone(nomor, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
		if gagal != nil { fmt.Println("GoError:",gagal); return }
		fmt.Println(nomor,">",linkingCode)
    } else {
        err = client.Connect()
        if err != nil { fmt.Println("GoError:",err); return }
        fmt.Println(nomor,">","Connected")
        client.SendPresence(types.PresenceAvailable)
        //client.SendPresence(types.PresenceUnavailable)
    }
}