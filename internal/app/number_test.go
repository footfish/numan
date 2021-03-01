package app_test

//Run tests from numan root with;
//go test -v -cover ./internal/app

import (
	"testing"
	"time"

	"github.com/footfish/numan"
	. "github.com/footfish/numan/internal/app"
)

// validPhoneNumbers
// country code (1 to 3 digits), no leading zero
// network code, with leading zero or '1' (2 to 5 digits)
// subscriber number (8 to 13 digits), no leading zero
var validPhoneNumbers = []numan.E164{
	{Cc: "1", Ndc: "01", Sn: "01111"},
	{Cc: "22", Ndc: "02", Sn: "22222222"},
	{Cc: "333", Ndc: "033", Sn: "3333333333"},
	{Cc: "1", Ndc: "0800", Sn: "123456789013"},
	{Cc: "1", Ndc: "1800", Sn: "123456789013"},
	{Cc: "1", Ndc: "01234", Sn: "123456789013"},
}
var invalidPhoneNumbers = []numan.E164{
	{Cc: "1", Ndc: "00", Sn: "11111111"},         //subsequent zero NDC
	{Cc: "22", Ndc: "02", Sn: "2222222A"},        //non numeric in SN number
	{Cc: "2A", Ndc: "02", Sn: "22222222"},        //non numeric in CC number
	{Cc: "22", Ndc: "0A", Sn: "22222222"},        //non numeric in NDC number
	{Cc: "0", Ndc: "033", Sn: "3333333333"},      //zero cc
	{Cc: "01", Ndc: "0800", Sn: "123456789013"},  //leading zero cc
	{Cc: "1", Ndc: "2123", Sn: "123456789013"},   //ndc not starting 0 or 1
	{Cc: "1234", Ndc: "033", Sn: "3333333333"},   //cc > 3
	{Cc: "1", Ndc: "123456", Sn: "123456789013"}, //ndc > 5
	{Cc: "1", Ndc: "01", Sn: "1234"},             //sn < 5
	{Cc: "1", Ndc: "01", Sn: "12345678901234"},   //sn > 13
	{Cc: "", Ndc: "033", Sn: "3333333333"},       //missing cc
	{Cc: "333", Ndc: "", Sn: "3333333333"},       //missing ndc
	{Cc: "333", Ndc: "033", Sn: ""},              //missing sn
}

func TestDelete(t *testing.T) {
	nu := HelperNewNumberService(t)
	defer nu.Close()
	//Add
	if err := nu.Add(&numan.Number{E164: validPhoneNumbers[0], Domain: "anydomain.com", Carrier: "anycarrier"}); err != nil {
		t.Fatal(err)
	}
	//Delete
	if err := nu.Delete(&validPhoneNumbers[0]); err != nil {
		t.Fatal(err)
	}
	if foundNumber, err := nu.List(&numan.NumberFilter{E164: validPhoneNumbers[0]}); err != nil {
		t.Fatal(err)
	} else if want, got := 0, len(foundNumber); want != got {
		t.Fatalf("Delete, found %v, want %v", got, want)
	}
}

func TestReserve(t *testing.T) {
	t.Run("OkReserveNumber", func(t *testing.T) {
		nu := HelperNewNumberService(t)
		defer nu.Close()
		//Add
		if err := nu.Add(&numan.Number{E164: validPhoneNumbers[0], Domain: "anydomain.com", Carrier: "anycarrier"}); err != nil {
			t.Fatal(err)
		}
		//Reserve
		untilTS := time.Now().Unix() + (60 * 15) //15mins
		userID := 99
		if err := nu.Reserve(&validPhoneNumbers[0], &userID, &untilTS); err != nil {
			t.Fatal(err)
		}
		//Read & check
		if storedNumber, err := nu.List(&numan.NumberFilter{E164: validPhoneNumbers[0]}); err != nil {
			t.Fatal(err)
		} else if want, got := untilTS, storedNumber[0].Reserved; want != got { //Cc
			t.Fatalf("Reserved got %v, want %v", got, want)
		}
	})

}
func TestAdd(t *testing.T) {
	//Verifiy basic Add
	t.Run("OkAddDeletePhoneNumber", func(t *testing.T) {
		nu := HelperNewNumberService(t)
		defer nu.Close()
		//Verify a number can be added and read back
		if err := nu.Add(&numan.Number{E164: validPhoneNumbers[0], Domain: "anydomain.com", Carrier: "anycarrier"}); err != nil {
			t.Fatal(err)
		}
		if storedNumber, err := nu.List(&numan.NumberFilter{E164: validPhoneNumbers[0]}); err != nil {
			t.Fatal(err)
		} else if want, got := validPhoneNumbers[0].Cc, storedNumber[0].E164.Cc; want != got { //Cc
			t.Fatalf("Cc got %v, want %v", got, want)
		} else if want, got := validPhoneNumbers[0].Ndc, storedNumber[0].E164.Ndc; want != got { //Ndc
			t.Fatalf("Ndc got %v, want %v", got, want)
		} else if want, got := validPhoneNumbers[0].Sn, storedNumber[0].E164.Sn; want != got { //Sn
			t.Fatalf("Sn got %v, want %v", got, want)
		} else if want, got := "anydomain.com", storedNumber[0].Domain; want != got { //Domain
			t.Fatalf("got %v, want %v", got, want)
		} else if want, got := "anycarrier", storedNumber[0].Carrier; want != got { //Carrier
			t.Fatalf("got %v, want %v", got, want)
		} else if want, got := 0, storedNumber[0].UserID; want != got { //UserID
			t.Fatalf("UserId got %v, want %v", got, want)
		} else if want, got := false, storedNumber[0].Used; want != got { //Used
			t.Fatalf("Used got %v, want %v", got, want)
		} else if want, got := int64(0), storedNumber[0].PortedIn; want != got { //PortedIn
			t.Fatalf("PortedIn got %v, want %v", got, want)
		} else if want, got := int64(0), storedNumber[0].PortedOut; want != got { //PortedOut
			t.Fatalf("PortedOut got %v, want %v", got, want)
		} else if want, got := int64(0), storedNumber[0].Allocated; want != got { //Allocated
			t.Fatalf("Allocated got %v, want %v", got, want)
		} else if want, got := int64(0), storedNumber[0].DeAllocated; want != got { //DeAllocated
			t.Fatalf("DeAllocated got %v, want %v", got, want)
		} else if want, got := int64(0), storedNumber[0].Reserved; want != got { //Reserved
			t.Fatalf("Reserved got %v, want %v", got, want)
		}

	})

	//Verifiy required parameters
	t.Run("ErrRequiredFields", func(t *testing.T) {
		nu := HelperNewNumberService(t)
		defer nu.Close()
		if err := nu.Add(&numan.Number{E164: validPhoneNumbers[0], Domain: "", Carrier: "anycarrier"}); err == nil {
			t.Fatal("Number added without domain")
		} else if want, got := "Carrier & domain required", err.Error(); want != got {
			t.Fatalf("Error '"+err.Error()+"' does not match '%v'", want)
		}
		if err := nu.Add(&numan.Number{E164: validPhoneNumbers[0], Domain: "anydomain", Carrier: ""}); err == nil {
			t.Fatal("Number added without carrier")
		} else if want, got := "Carrier & domain required", err.Error(); want != got {
			t.Fatalf("Error '"+err.Error()+"' does not include '%v'", want)
		}
	})

	//Verify all valid phone numbers can be added
	t.Run("OkValidPhoneNumbers", func(t *testing.T) {
		nu := HelperNewNumberService(t)
		defer nu.Close()
		for _, phoneNumber := range validPhoneNumbers {
			if err := nu.Add(&numan.Number{E164: phoneNumber, Domain: "anydomain.com", Carrier: "anycarrier", Used: true}); err != nil {
				t.Fatalf(err.Error()+"number %v-%v-%v", phoneNumber.Cc, phoneNumber.Ndc, phoneNumber.Sn)
			}
		}
	})

	//Verify invalid phone numbers can't be added
	t.Run("ErrInvalidPhoneNumbers", func(t *testing.T) {
		nu := HelperNewNumberService(t)
		defer nu.Close()
		for _, phoneNumber := range invalidPhoneNumbers {
			if err := nu.Add(&numan.Number{E164: phoneNumber, Domain: "anydomain.com", Carrier: "anycarrier", Used: true}); err == nil {
				t.Fatalf("Added invalid number %v-%v-%v", phoneNumber.Cc, phoneNumber.Ndc, phoneNumber.Sn)
			}
		}

	})
}

// NewNumberService instantiates a new NuService.
func HelperNewNumberService(t *testing.T) *NumberService {
	t.Helper()
	return NewNumberService(":memory:")
}
