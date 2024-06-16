package wa
import (
    "fmt"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/store/sqlstore"
    "go.mau.fi/whatsmeow/types"
    waLog "go.mau.fi/whatsmeow/util/log"
    _ "github.com/mattn/go-sqlite3"
)

var Conns = make(map[string]*whatsmeow.Client)
func Connect(nomor string, eventHandler func(evt interface{})) {
    dbLog := waLog.Stdout("Database", "ERROR", true)
    container, err := sqlstore.New("sqlite3", "file:"+nomor+".db?_foreign_keys=on", dbLog)
    if err != nil { fmt.Println("GoError:",err); return }
    deviceStore, err := container.GetFirstDevice()
    if err != nil { fmt.Println("GoError:",err); return }
    clientLog := waLog.Stdout("Client", "ERROR", true)
    client := whatsmeow.NewClient(deviceStore, clientLog)
    Conns[nomor] = client
    client.AddEventHandler(eventHandler)
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