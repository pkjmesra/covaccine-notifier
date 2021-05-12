package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	
)
type BookingSlot struct {
		Available 	bool
		Preferred 	bool
		CenterID 	int
		CenterName  string
		SessionID 	string
		Slot 		string
		Description string
	}

var (
	pinCode, state, district, email, password, notificationFile, whatsAppRemoteNum, date string
	slotsAvailable bool
	age, interval, bookingCenterId int
	bookingSlot *BookingSlot

	rootCmd = &cobra.Command{
		Use:   "covaccine-notifier [FLAGS]",
		Short: "CoWIN Vaccine availability notifier India",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(args)
		},
	}
)

const (
	pinCodeEnv        		= "PIN_CODE"
	stateNameEnv      		= "STATE_NAME"
	districtNameEnv   		= "DISTRICT_NAME"
	ageEnv            		= "AGE"
	emailIDEnv        		= "EMAIL_ID"
	emailPasswordEnv  		= "EMAIL_PASSOWORD"
	searchIntervalEnv 		= "SEARCH_INTERVAL"
	notificationMP3FileEnv 	= "NOTIFICATION_MP3_FILE"
	whatsappRemoteNumEnv 	= "REMOTE_WHATSAPP_NUM"
	bookingCenterIdEnv 		= "BOOKING_CENTERID"

	defaultSearchInterval = 60
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&pinCode, "pincode", "c", os.Getenv(pinCodeEnv), "Search by pin code")
	rootCmd.PersistentFlags().StringVarP(&state, "state", "s", os.Getenv(stateNameEnv), "Search by state name")
	rootCmd.PersistentFlags().StringVarP(&district, "district", "d", os.Getenv(districtNameEnv), "Search by district name")
	rootCmd.PersistentFlags().IntVarP(&age, "age", "a", getIntEnv(ageEnv), "Search appointment for age")
	rootCmd.PersistentFlags().StringVarP(&email, "email", "e", os.Getenv(emailIDEnv), "Email address to send notifications")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", os.Getenv(emailPasswordEnv), "Email ID password for auth")
	rootCmd.PersistentFlags().IntVarP(&interval, "interval", "i", getIntEnv(searchIntervalEnv), fmt.Sprintf("Interval to repeat the search. Default: (%v) second", defaultSearchInterval))
	rootCmd.PersistentFlags().StringVarP(&notificationFile, "notificationFile", "n", os.Getenv(notificationMP3FileEnv), "Specify a local MP3 file to play when a slot is available")
	rootCmd.PersistentFlags().StringVarP(&whatsAppRemoteNum, "whatsAppRemoteNum", "w", os.Getenv(whatsappRemoteNumEnv), "Specify a remote WhatsApp mobile number")
	rootCmd.PersistentFlags().IntVarP(&bookingCenterId, "bookingCenterId", "b", getIntEnv(bookingCenterIdEnv), "Preferred booking center Id")

}

// Execute executes the main command
func Execute() error {
	slotsAvailable = false
	bookingSlot := BookingSlot{}
	bookingSlot.Available = slotsAvailable
	return rootCmd.Execute()
}

func checkFlags() error {
	if len(pinCode) == 0 &&
		len(state) == 0 &&
		len(district) == 0 {
		return errors.New("Please pass one of the pinCode or state & district name combination options")
	}
	if len(pinCode) == 0 && (len(state) == 0 || len(district) == 0) {
		return errors.New("Missing state or district name option")
	}
	if age == 0 {
		return errors.New("Missing age option")
	}
	if len(email) == 0 || len(password) == 0 {
		return errors.New("Missing email creds")
	}
	if interval == 0 {
		interval = defaultSearchInterval
	}
	return nil
}

func main() {
	Execute()
}

func getIntEnv(envVar string) int {
	v := os.Getenv(envVar)
	if len(v) == 0 {
		return 0
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func Run(args []string) error {
	if err := checkFlags(); err != nil {
		return err
	}
	if err := checkSlots(); err != nil {
		return err
	}
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := checkSlots(); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkSlots() error {
	// Search for slots
	var err error
	var bk *BookingSlot
	if len(pinCode) != 0 {
		bk, err = searchByPincode(pinCode)
		return err
	}
	bk, err = searchByStateDistrict(age, state, district)
	if bk.Available && bk.Preferred {
		if err = bookSlot(bk); err != nil {
			return err
		}
	}
	return err
}

func bookSlot(availableSlot *BookingSlot) error {
	log.Printf("Going to book a slot for CenterID :%d, SessionID: %s, Slot: %s", availableSlot.CenterID, availableSlot.SessionID, availableSlot.Slot)
	return nil
}