package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
	"strings"

	// "github.com/Rhymen/go-whatsapp/binary/proto"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

type waHandler struct {
	c *whatsapp.Conn
}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (h *waHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		// log.Printf("error occoured: %v\n", err)
	}
}

//Optional to be implemented. Implement HandleXXXMessage for the types you need.
func (*waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	numbers := strings.Split(whatsAppRemoteNum, ",")
	oneFromRegisteredNumber := false
	for _, remoteNumber := range numbers {
		if message.Info.RemoteJid == remoteNumber + "@s.whatsapp.net"{
			oneFromRegisteredNumber = true
		}
	}
	shouldProcess := oneFromRegisteredNumber // || message.Info.FromMe
	if  shouldProcess {
		// fmt.Printf("%v %v %v %v\n\t%v\n", message.Info.Timestamp, message.Info.Id, message.Info.RemoteJid, message.ContextInfo.QuotedMessageID, message.Text)
		// if len(message.Text) == 36 {
		// 	// Must be txnid.
		// 	txnIdReceived(message.Text)
		// }
		// Parse only OTP
		if len(message.Text) == 6 {
			// Must be OTP. Let's confirm the OTP
			OTPReceived(message.Text)
		}
		// if len(message.Text) >= 400 {
		// 	// Must be bearer token.
		// 	bearerTokenReceived(message.Text)
		// }
	}
}

func sendwhatsapptext(textMessage string) {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	//Add handler
	wac.AddHandler(&waHandler{wac})

	err = login(wac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		return
	}

	<-time.After(10 * time.Second)

	// previousMessage := "😘"
	// quotedMessage := proto.Message{
	// 	Conversation: &previousMessage,
	// }

	// ContextInfo := whatsapp.ContextInfo{
	// 	QuotedMessage:   &quotedMessage,
	// 	QuotedMessageID: "",
	// 	Participant:     "", //Who sent the original message
	// }

	numbers := strings.Split(whatsAppRemoteNum, ",")

	for _, remoteNumber := range numbers {
		msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: remoteNumber + "@s.whatsapp.net",
			FromMe: true,
		},
		// ContextInfo: ContextInfo,
		Text: textMessage,
		}

		msgId, err := wac.Send(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error sending message: %v", err)
			// os.Exit(1)
		} else {
			fmt.Println("Message Sent -> ID : " + msgId)
		}	
	}
}

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed from " + os.TempDir() + "whatsappSession.gob"  + ": %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
		fmt.Println("Login done -> ID : " + string(session.Wid))
	}

	loggedInWANumber = strings.TrimSuffix(string(session.Wid), "@c.us")
	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}