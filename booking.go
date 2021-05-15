package main

import (
	"encoding/json"
	"log"

	// "github.com/pkg/errors"
)

type AppointmentConfirmation struct {
	ConfirmationId string `json:"appointment_confirmation_no"`
	// Message string `json:"message"`
}

type ApptRequestError struct {
	Error    	string `json:"error"`
	ErrorCode  	string `json:"errorCode"`
}

func bookAppointment(beneficiaryList *BeneficiaryList) {
	beneficiaryId := ""
	outer:
	for _, beneficiary := range beneficiaryList.Beneficiaries {
		if beneficiary.VaccinationStat == "Not Vaccinated" {
			beneficiaryId = beneficiary.ReferenceID
			break outer
		}
	}
	beneficiaries:= make([]string,1)
	beneficiaries[0] = beneficiaryId
	postBody := map[string]interface{}{"center_id": bookingSlot.CenterID, "dose": 1, "captcha": "","session_id": bookingSlot.SessionID, "slot": bookingSlot.Slot, "beneficiaries": beneficiaries}
	bodyBytes, err := queryServer(scheduleURLFormat, "POST", postBody)
	cnf := AppointmentConfirmation{}
	if err = json.Unmarshal(bodyBytes, &cnf); err != nil {
		aptErr := ApptRequestError{}
		if err = json.Unmarshal(bodyBytes, &aptErr); err != nil {
			log.Printf("Error scheduling: %s", err.Error())
		}
		log.Printf("ErrorCode:%s , Error:%s", aptErr.ErrorCode, aptErr.Error)
		return
	}
	log.Printf("AppointmentID confirmed:%s", cnf.ConfirmationId)
}
